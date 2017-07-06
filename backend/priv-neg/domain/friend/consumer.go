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
		true,         // auto-ack
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

func (c *Consumer) process(d amqp.Delivery) {
	start := time.Now()

	dbUser := &domain.DBUser{}
	json.Unmarshal(d.Body, dbUser)

	dbUser, err := c.userDB.FindByID(dbUser.ID)
	if err != nil {
		return
	}

	c.logger.Printf("Started processing cliques for %s", dbUser.ID)

	// k=3
	// 1. Get all existing cliques for User
	// 2. Iterate through the users involved in those cliques and discard from new clique detection.
	// 3. With reduced users, look for mutual friends and form new cliques.
	// 4. If new clique contains users in existing clique, get that clique and add reduced user to it.
	// 5. remove users in new clique from reduced users

	existingFriendsInCliques := []string{}
	for _, userClique := range dbUser.DBUserCliques {
		clique, err := c.cliqueDB.FindCliqueByID(userClique.CliqueID)
		if err != nil {
			c.logger.Printf("error: %v", err)
		}
		existingFriendsInCliques = append(existingFriendsInCliques, clique.GetUserIDs()...)
	}

	c.logger.Printf("debug: Got %d friends already in cliques", len(existingFriendsInCliques))

	allFriendIDs := c.getFacebookFriendsForUser(dbUser, "", nil)

	reducedFriendIDs := []string{}

	for _, friendID := range allFriendIDs {
		if !isIn(friendID, existingFriendsInCliques) {
			reducedFriendIDs = append(reducedFriendIDs, friendID)
		}
		cacheUser := domain.CacheUserFromDatabaseUser(dbUser)
		cacheFriendship := &domain.CacheFriendship{ID: friendID}
		c.friendRedis.Save(cacheUser, cacheFriendship)

		// Add to tie-strength queue
		queueFriendship := &domain.QueueFriendship{
			From: dbUser.ID,
			To:   friendID,
		}
		c.tieStrengthPublisher.Publish(queueFriendship)
	}

	alreadyUserReducedUsers := []string{}

	for _, friendID := range reducedFriendIDs {
		if isIn(friendID, alreadyUserReducedUsers) {
			c.logger.Printf("debug: Skipping %s", friendID)
			break // skip user
		}
		friendFriends := c.friendRedis.GetFriendIDsForAUserID(friendID)

		mutualFriends := arrayUnion(allFriendIDs, friendFriends)
		c.logger.Printf("debug: Got mutual friends: %v", mutualFriends)

		if len(mutualFriends) >= 1 {
			if isIn(mutualFriends[0], reducedFriendIDs) {
				c.logger.Printf("Found a new clique size %d for %s", len(mutualFriends)+2, dbUser.ID)
				// If completely new clique
				clique := domain.NewCacheClique()
				dbClique := domain.DBCliqueFromCacheClique(clique)
				c.cliqueDB.Save(dbClique)

				// Add clique to all members
				cliqueMembers := append(append(mutualFriends, dbUser.ID), friendID)
				for _, userID := range cliqueMembers {
					c.friendRedis.AddCliqueToUserID(userID, clique)
					dbUserClique := domain.DBUserCliqueFromCacheCliqueAndUserID(clique, userID)
					c.cliqueDB.SaveUserClique(dbUserClique)
					alreadyUserReducedUsers = append(alreadyUserReducedUsers, userID)
					webClique := domain.WSClique{
						ID:      dbUserClique.CliqueID,
						Name:    "",
						UserIDs: cliqueMembers,
					}
					c.userRedis.Publish(userID, "clique", webClique)
				}

			} else {
				c.logger.Printf("Found an existing clique size %d for %s", len(mutualFriends)+1, dbUser.ID)
				// Add to existing clique
				// Get set of cliques for each user from cache and compare.
				userCliques := map[int][]string{} // i: []cliqueID

				mutualCliqueIDs := []string{}

				for i, friendID := range mutualFriends {
					cliqueIDs := c.friendRedis.GetCliqueIDsForAUserID(friendID)
					if len(cliqueIDs) == 0 {
						break
					}
					userCliques[i] = cliqueIDs
					if i == 0 {
						mutualCliqueIDs = cliqueIDs
					} else if i < len(mutualFriends) {
						mutualCliqueIDs = arrayUnion(mutualCliqueIDs, cliqueIDs)
					}
				}

				// TODO: Merge cliques when necessary
				// Merging cliques:
				// Create a new clique.
				// Go through all of the mutual cliques, add all the users and categories.
				// Remove old cliques
				c.logger.Printf("warning: This should == 1: %v (if it's greater, then the cliques need merging)", mutualCliqueIDs)

				newClique := domain.NewCacheClique()
				dbClique := domain.DBCliqueFromCacheClique(newClique)
				c.cliqueDB.Save(dbClique)
				userCliquesCategories := map[string][]domain.DBCategory{}
				userCliquesNames := map[string]string{}

				for _, mutualCliqueID := range mutualCliqueIDs {
					mutualDBClique, err := c.cliqueDB.FindCliqueByID(mutualCliqueID)
					if err != nil {
						break
					}
					for _, userClique := range mutualDBClique.DBUserCliques {
						if _, ok := userCliquesCategories[userClique.UserID]; !ok {
							userCliquesCategories[userClique.UserID] = []domain.DBCategory{}
						}
						if _, ok := userCliquesNames[userClique.UserID]; !ok {
							userCliquesNames[userClique.UserID] = ""
						}

						userCliquesCategories[userClique.UserID] = append(userCliquesCategories[userClique.UserID], userClique.Categories...)
						userCliquesNames[userClique.UserID] += fmt.Sprintf(" %s", userClique.Name)
						c.friendRedis.RemoveCliqueByIDFromUserID(userClique.UserID, userClique.CliqueID)
					}
					c.cliqueDB.DeleteCliqueByID(mutualCliqueID)
				}

				for userID, dbCategories := range userCliquesCategories {
					dbUserClique := domain.DBUserClique{
						CliqueID:   newClique.ID,
						UserID:     userID,
						Categories: dbCategories,
						Name:       userCliquesNames[userID],
					}
					c.cliqueDB.SaveUserClique(&dbUserClique)
					userCacheClique := domain.CacheClique{
						ID:             newClique.ID,
						Name:           userCliquesNames[userID],
						Categories:     []string{},
						UserCategories: []string{},
					}
					for _, cat := range dbCategories {
						if cat.UserID == "none" {
							userCacheClique.Categories = append(userCacheClique.Categories, cat.Name)
						} else {
							userCacheClique.UserCategories = append(userCacheClique.UserCategories, cat.Name)
						}
					}
					c.friendRedis.AddCliqueToUserID(userID, &userCacheClique)
					c.userRedis.Publish(userID, "clique", userCacheClique)
				}

				mutualCliqueID := newClique.ID
				clique := &domain.CacheClique{
					ID:             mutualCliqueID,
					Name:           "",
					Categories:     []string{},
					UserCategories: []string{},
				}

				c.friendRedis.AddCliqueToUserID(friendID, clique)
				dbUserClique := domain.DBUserCliqueFromCacheCliqueAndUserID(clique, friendID)
				c.cliqueDB.SaveUserClique(dbUserClique)

				// Update WS
				c.userRedis.Publish(friendID, "clique", clique)
				//dbClique, _ := c.cliqueDB.FindCliqueByID(newClique.ID)
				//for _, userID := range dbClique.GetUserIDs() {
				//	c.userRedis.Publish(userID, "clique", clique)
				//}
			}

		}

	}

	elapsed := time.Since(start)
	c.logger.Printf("Processed cliques for %s in %s", dbUser.ID, elapsed)
}

func arrayUnion(a []string, b []string) []string {
	c := []string{}

	for _, aV := range a {
		for _, bV := range b {
			if aV == bV && !isIn(aV, c) {
				c = append(c, aV)
			}
		}
	}

	return c
}

func isIn(needle string, haystack []string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
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
