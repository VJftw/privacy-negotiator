package friend

import (
	"log"
	"time"

	"encoding/json"

	"fmt"
	"net/http"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/utils"
	"github.com/streadway/amqp"
)

// Consumer - Consumer for getting Community Detection.
type Consumer struct {
	queue                amqp.Queue
	channel              *amqp.Channel
	logger               *log.Logger
	userDB               *user.DBManager
	userRedis            *user.RedisManager
	friendRedis          *RedisManager
	cliqueDB             *DBManager
	tieStrengthPublisher *TieStrengthPublisher
}

// NewConsumer - Returns a new consumer.
func NewConsumer(
	queueLogger *log.Logger,
	ch *amqp.Channel,
	userDBManager *user.DBManager,
	userRedisManager *user.RedisManager,
	friendRedisManager *RedisManager,
	cliqueDBManager *DBManager,
	tieStrengthPublisher *TieStrengthPublisher,
) *Consumer {
	queue, err := ch.QueueDeclare(
		"community-detection", // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")
	return &Consumer{
		logger:               queueLogger,
		channel:              ch,
		queue:                queue,
		userDB:               userDBManager,
		userRedis:            userRedisManager,
		friendRedis:          friendRedisManager,
		cliqueDB:             cliqueDBManager,
		tieStrengthPublisher: tieStrengthPublisher,
	}
}

// Consume - Processes items from the Queue.
func (c *Consumer) Consume() {
	msgs, err := c.channel.Consume(
		c.queue.Name, // queue
		"",           // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	utils.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			c.process(d)
		}
	}()

	c.logger.Printf("Worker: %s waiting for messages. To exit press CTRL+C", c.queue.Name)
	<-forever
}

// 1. Find new friends. Save bidirectional relationships with them.
// 2. Build local graph. Max 1 step (immediate friends as we're only interested in finding cliques the user can be part of)
// 3. For each user, find mutual friends and form a clique (reduce as you go along so dupe cliques aren't formed)
//
// 4. Compare these cliques to existing ones the user is in.
// 		if they exist, ignore them.
// 		if the existing ones are a smaller subset of a new one. Migrate to new and delete old.
//      else form a new clique.
//

// Process - Processes queue
func (c *Consumer) Process(d amqp.Delivery) {
	start := time.Now()

	dbUser := &domain.DBUser{}
	json.Unmarshal(d.Body, dbUser)
	dbUser, err := c.userDB.FindByID(dbUser.ID)
	if err != nil {
		return
	}

	c.logger.Printf("Processing cliques for %s", dbUser.ID)

	// 1.
	// Find new friends and save bidirectional relationships in cache with them.
	currentFriendIDs := c.friendRedis.GetFriendIDsForAUserID(dbUser.ID)
	c.logger.Printf("Got $d friends for %s from Cache", len(currentFriendIDs), dbUser.ID)
	allFriendIDs := c.getFacebookFriendsForUser(dbUser, "", nil)
	c.logger.Printf("Got $d friends for %s from Graph API", len(allFriendIDs), dbUser.ID)
	for _, fID := range allFriendIDs {
		if !utils.IsIn(fID, currentFriendIDs) {
			// Save new bidirectional friendship
			uFriendship := &domain.CacheFriendship{ID: fID}
			c.friendRedis.Save(dbUser.ID, uFriendship)
			fFriendship := &domain.CacheFriendship{ID: dbUser.ID}
			c.friendRedis.Save(fID, fFriendship)

			// Add to current friends as they are now saved
			currentFriendIDs = append(currentFriendIDs, fID)
		}
	}

	// 2.
	// Build local graph
	c.logger.Printf("info: building local friend graph")
	localFriendGraph := map[string][]string{} // <userID: []friendIDs + userID>
	uFriendList := append(currentFriendIDs, dbUser.ID)
	for _, fID := range currentFriendIDs {
		localFriendGraph[fID] = append(c.friendRedis.GetFriendIDsForAUserID(fID), fID)
	}
	c.logger.Printf("debug: built local friend graph: %v", localFriendGraph)

	// 3.
	// Get existing cliques for next step
	userCliqueIDs := c.friendRedis.GetCliqueIDsForAUserID(dbUser.ID)
	existingDBCliques := []*domain.DBClique{}
	for _, userCliqueID := range userCliqueIDs {
		dbClique, _ := c.cliqueDB.FindCliqueByID(userCliqueID)
		existingDBCliques = append(existingDBCliques, dbClique)
	}

	// 4.
	// Compare user's friends to each friend's friends using array union. If result is >= 3, we can form a clique.
	// If a clique is found. Compare to existing ones:
	//  	- If a smaller subset already exists, migrate to the new one.
	//		- If an identical one already exists, ignore the new one.
	// 		- else form a new clique normally.
	inClique := []string{}
	for fID, fFriendList := range localFriendGraph {
		if utils.IsIn(fID, inClique) {
			break
		}
		c.logger.Printf("debug: comparing %s: %v to %s: %v", dbUser.ID, uFriendList, fID, fFriendList)
		mutualFriends := utils.ArrayUnion(uFriendList, fFriendList)
		if len(mutualFriends) >= 3 {
			c.logger.Printf("debug: found potential clique: %v", mutualFriends)
			if isValidClique(mutualFriends, dbUser.ID, localFriendGraph) {
				c.logger.Printf("debug: found new clique: %v", mutualFriends)
				existingDBCliques = c.processClique(mutualFriends, existingDBCliques)
				// Now we can ignore friend keys that are in this clique
				inClique = append(inClique, mutualFriends...)
			}
		}
	}

	c.logger.Printf("Processed cliques for %s in %s", dbUser.ID, time.Since(start))
	d.Ack(false)
}

func isValidClique(newClique []string, userID string, localFriendGraph map[string][]string) bool {
	for _, graphKey := range newClique {
		if graphKey == userID {
			break
		}
		for _, uID := range newClique {
			if !utils.IsIn(uID, localFriendGraph[graphKey]) {
				return false
			}
		}
	}

	return true
}

func (c *Consumer) processClique(newClique []string, existingCliques []*domain.DBClique) []*domain.DBClique {
	migrateUserCliquesCategories := map[string][]domain.DBCategory{} // userID: categories
	migrateUserCliquesNames := map[string]string{}                   // userID: names

	returnCliques := []*domain.DBClique{}
	c.logger.Printf("EXISTING CLIQUES: %v", existingCliques)
	var dbClique *domain.DBClique
	for _, existingClique := range existingCliques {
		c.logger.Printf("EXISTING CLIQUE: %s: %v", existingClique.ID, existingClique.GetUserIDs())

		if utils.IsSubset(existingClique.GetUserIDs(), newClique) && len(newClique) == len(existingClique.GetUserIDs()) {
			// Cliques are identical. Do nothing.
			c.logger.Printf("debug: new clique is identical to %s", existingClique.ID)
			return existingCliques
		} else if utils.IsSubset(newClique, existingClique.GetUserIDs()) {
			// New clique is a subset
			c.logger.Printf("debug: new clique is a subset to %s", existingClique.ID)
			return existingCliques
		} else if utils.IsSubset(existingClique.GetUserIDs(), newClique) && len(newClique) > len(existingClique.GetUserIDs()) {
			// existing clique to migration for new clique and remove old clique.
			for _, userClique := range existingClique.DBUserCliques {
				c.logger.Printf("debug: copying user clique %s:%s for migration", userClique.UserID, userClique.CliqueID)
				if _, ok := migrateUserCliquesCategories[userClique.UserID]; !ok {
					migrateUserCliquesCategories[userClique.UserID] = []domain.DBCategory{}
				}
				if _, ok := migrateUserCliquesNames[userClique.UserID]; !ok {
					migrateUserCliquesNames[userClique.UserID] = ""
				}
				migrateUserCliquesCategories[userClique.UserID] = append(migrateUserCliquesCategories[userClique.UserID], userClique.Categories...)
				migrateUserCliquesNames[userClique.UserID] += fmt.Sprintf(" %s", userClique.Name)
				c.friendRedis.RemoveCliqueByIDFromUserID(userClique.UserID, userClique.CliqueID)
				c.cliqueDB.DeleteUserClique(userClique.CliqueID, userClique.UserID)
			}
			c.logger.Printf("info: deleting smaller subset clique: %s", existingClique.ID)
			c.cliqueDB.DeleteCliqueByID(existingClique.ID)
		} else {
			returnCliques = append(returnCliques, existingClique)
		}
	}

	// form new clique
	dbClique = domain.NewDBClique()

	newCacheClique := &domain.CacheClique{
		ID: dbClique.ID,
	}
	c.cliqueDB.Save(dbClique)
	for _, uID := range newClique {
		dbUserClique := domain.DBUserCliqueFromCacheCliqueAndUserID(newCacheClique, uID)
		if _, ok := migrateUserCliquesCategories[uID]; ok {
			// if migrating categories
			dbUserClique.Categories = append(dbUserClique.Categories, migrateUserCliquesCategories[uID]...)
		}
		if _, ok := migrateUserCliquesNames[uID]; ok {
			// if migrating names
			dbUserClique.Name += migrateUserCliquesNames[uID]
		}

		cacheClique, _ := domain.CacheCliqueFromDBUserClique(dbUserClique)
		c.friendRedis.AddCliqueToUserID(uID, cacheClique)
		c.cliqueDB.SaveUserClique(dbUserClique)
		dbClique.DBUserCliques = append(dbClique.DBUserCliques, *dbUserClique)
		wsClique := domain.WSClique{
			ID:      dbUserClique.CliqueID,
			Name:    "",
			UserIDs: newClique,
		}
		c.userRedis.Publish(uID, "clique", wsClique)
	}

	return append(returnCliques, dbClique)
}

func (c *Consumer) process(d amqp.Delivery) {
	c.Process(d)
}

func (c *Consumer) getFacebookFriendsForUser(user *domain.DBUser, offset string, friendIds []string) []string {

	if friendIds == nil {
		friendIds = []string{}
	}

	res, _ := http.Get(fmt.Sprintf(
		"https://graph.facebook.com/v2.9/%s/friends?access_token=%s&limit=500%s",
		user.ID,
		user.LongLivedToken,
		offset,
	))

	responseFriends := &responseFriends{}
	err := json.NewDecoder(res.Body).Decode(responseFriends)
	defer res.Body.Close()
	if err != nil {
		log.Printf("Error: %s", err)
	}

	for _, responseFriend := range responseFriends.Data {
		friendIds = append(friendIds, responseFriend.ID)
	}

	if responseFriends.Paging.Cursors.After != "" {
		return c.getFacebookFriendsForUser(
			user,
			fmt.Sprintf("&after=%s", responseFriends.Paging.Cursors.After),
			friendIds,
		)
	}

	return friendIds
}

type responseFriends struct {
	Data   []responseFriend `json:"data"`
	Paging responsePaging   `json:"paging"`
}

type responseFriend struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type responsePaging struct {
	Cursors responseCursors `json:"cursors"`
}

type responseCursors struct {
	Before string `json:"before"`
	After  string `json:"after"`
}
