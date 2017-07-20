package photo

import (
	"encoding/json"
	"log"
	"time"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/friend"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/utils"
	"github.com/streadway/amqp"
)

// ConflictConsumer - ConflictConsumer for conflict detection and suggested resolution.
type ConflictConsumer struct {
	queue       amqp.Queue
	channel     *amqp.Channel
	logger      *log.Logger
	friendDB    *friend.DBManager
	userDB      *user.DBManager
	photoDB     *DBManager
	photoRedis  *RedisManager
	userRedis   *user.RedisManager
	friendRedis *friend.RedisManager
}

// NewConflictConsumer - Returns a new consumer.
func NewConflictConsumer(
	queueLogger *log.Logger,
	ch *amqp.Channel,
	friendDBManager *friend.DBManager,
	userDBManager *user.DBManager,
	photoDBManager *DBManager,
	photoRedisManager *RedisManager,
	userRedisManager *user.RedisManager,
	friendRedisManager *friend.RedisManager,
) *ConflictConsumer {
	queue, err := ch.QueueDeclare(
		"conflict-detection-and-resolution", // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")
	return &ConflictConsumer{
		logger:      queueLogger,
		channel:     ch,
		queue:       queue,
		friendDB:    friendDBManager,
		userDB:      userDBManager,
		photoDB:     photoDBManager,
		photoRedis:  photoRedisManager,
		userRedis:   userRedisManager,
		friendRedis: friendRedisManager,
	}
}

// Consume - Processes items from the Queue.
func (c *ConflictConsumer) Consume() {
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

func (c *ConflictConsumer) process(d amqp.Delivery) {
	start := time.Now()

	dbPhoto := domain.DBPhoto{}
	json.Unmarshal(d.Body, &dbPhoto)

	taggedUserAllowedUsers := map[string][]string{} // taggedUserID: cliqueUserID
	taggedUserBlockedUsers := map[string][]string{} // taggedUserID: cliqueUserID

	c.logger.Printf("debug: TaggedUsers: %d, Categories: %d", len(dbPhoto.TaggedUsers), len(dbPhoto.Categories))
	taggedUserIDs := []string{}
	for _, taggedUser := range dbPhoto.TaggedUsers {
		taggedUserIDs = append(taggedUserIDs, taggedUser.ID)
	}

	for _, taggedUser := range dbPhoto.TaggedUsers {
		userInitialAllowed := []string{taggedUser.ID}
		userInitialBlocked := []string{}
		taggedUserCliques, _ := c.friendDB.GetUserCliquesByUser(taggedUser)
		for _, taggedUserClique := range taggedUserCliques {
			c.logger.Printf("debug: UserCliqueCategories: %v", taggedUserClique.Categories)
			hasCategory := false
			for _, cliqueCat := range taggedUserClique.Categories {
				c.logger.Printf("debug: Allowing UserClique: %s", taggedUserClique.CliqueID)
				if isIn(cliqueCat, dbPhoto.Categories) {
					hasCategory = true
					break
				}
			}
			clique, _ := c.friendDB.FindCliqueByID(taggedUserClique.CliqueID)
			if hasCategory {
				userInitialAllowed = append(userInitialAllowed, clique.GetUserIDs()...)
			} else {
				userInitialBlocked = append(userInitialBlocked, clique.GetUserIDs()...)
			}
		}

		taggedUserAllowedUsers[taggedUser.ID] = []string{}
		taggedUserBlockedUsers[taggedUser.ID] = []string{}

		for _, blockedUserID := range userInitialBlocked {
			if !utils.IsIn(blockedUserID, userInitialAllowed) && !utils.IsIn(blockedUserID, taggedUserIDs) {
				taggedUserBlockedUsers[taggedUser.ID] = append(taggedUserBlockedUsers[taggedUser.ID], blockedUserID)
			}
		}
		for _, allowedUserID := range userInitialAllowed {
			taggedUserAllowedUsers[taggedUser.ID] = append(taggedUserAllowedUsers[taggedUser.ID], allowedUserID)

		}
	}

	c.logger.Printf("debug: Allowed Users %v", taggedUserAllowedUsers)
	c.logger.Printf("debug: Blocked Users %v", taggedUserBlockedUsers)

	// Find Conflicts
	cachePhoto, _ := c.photoRedis.FindByID(dbPhoto.ID)
	cachePhoto.AllowedUserIDs = []string{}
	cachePhoto.BlockedUserIDs = []string{}
	conflicts := map[string]domain.DBConflict{} //targetID: dbConflict.
	// Update AllowedUsers for CachedPhoto as well as finding conflicts.
	for taggedUserIDAllowed, allowedUserIDs := range taggedUserAllowedUsers {
		for _, allowedUserID := range allowedUserIDs {
			conflict := false
			for taggedUserIDBlocked, blockedUserIDs := range taggedUserBlockedUsers {
				for _, blockedUserID := range blockedUserIDs {
					if allowedUserID == blockedUserID {
						// conflict
						c.logger.Printf("Found conflict between %s and %s with %s", taggedUserIDAllowed, taggedUserIDBlocked, allowedUserID)
						var dbConflict domain.DBConflict

						if val, ok := conflicts[allowedUserID]; ok {
							// target already conflicted
							dbConflict = val
							if !isInUsers(taggedUserIDAllowed, dbConflict.Parties) {
								taggedUserAllowed, _ := c.userDB.FindByID(taggedUserIDAllowed)
								dbConflict.Parties = append(dbConflict.Parties, *taggedUserAllowed)
							}
							if !isInUsers(taggedUserIDBlocked, dbConflict.Parties) {
								taggedUserBlocked, _ := c.userDB.FindByID(taggedUserIDBlocked)
								dbConflict.Parties = append(dbConflict.Parties, *taggedUserBlocked)
							}
						} else {
							dbConflict = domain.NewDBConflict()
							allowedUser, _ := c.userDB.FindByID(allowedUserID)
							dbConflict.Target = *allowedUser
							dbConflict.Photo = dbPhoto
							dbConflict.PhotoID = dbPhoto.ID
							taggedUserAllowed, _ := c.userDB.FindByID(taggedUserIDAllowed)
							taggedUserBlocked, _ := c.userDB.FindByID(taggedUserIDBlocked)
							dbConflict.Parties = append(dbConflict.Parties, *taggedUserBlocked, *taggedUserAllowed)
						}

						conflicts[dbConflict.Target.ID] = dbConflict
						conflict = true
					}
				}
			}
			if !conflict && !utils.IsIn(allowedUserID, cachePhoto.AllowedUserIDs) {
				cachePhoto.AllowedUserIDs = append(cachePhoto.AllowedUserIDs, allowedUserID)
			}
		}
	}

	// Update BlockedUsers
	for _, blockedUserIDs := range taggedUserBlockedUsers {
		for _, blockedUserID := range blockedUserIDs {
			if _, ok := conflicts[blockedUserID]; !ok && !utils.IsIn(blockedUserID, cachePhoto.BlockedUserIDs) {
				cachePhoto.BlockedUserIDs = append(cachePhoto.BlockedUserIDs, blockedUserID)
			}
		}
	}
	cachePhoto.Conflicts = []domain.CacheConflict{}
	if len(conflicts) > 0 {

		for _, dbConflict := range conflicts {
			// We have conflicts. Save to Cache and DB and suggest resolution
			cacheConflict := domain.CacheConflictFromDBConflict(dbConflict)
			err := c.photoDB.SaveConflict(&dbConflict)
			if err != nil {
				c.logger.Printf("error: %v", err)
			}
			// Add party users with all sentiment to the user, their tie-strength is the weight of their vote.
			// Positive outcome means they should be allowed. Negative means they should be blocked.
			// (I love Democracy. I love the Republic)
			pointsAllow := 0
			for _, partyUser := range dbConflict.Parties {
				cacheUser := domain.CacheUserFromDatabaseUser(&partyUser)
				cacheFriendship, err := c.friendRedis.FindByIDAndUser(cacheConflict.Target, cacheUser)
				if err != nil {
					c.logger.Printf("debug: %s:%s do not have a friendship", partyUser.ID, cacheConflict.Target)
					break
				}
				if shouldAllow(taggedUserAllowedUsers, taggedUserBlockedUsers, partyUser.ID, cacheConflict.Target) {
					pointsAllow += cacheFriendship.TieStrength
					cacheConflict.Reasoning = append(
						cacheConflict.Reasoning,
						domain.Reason{
							UserID: partyUser.ID,
							Vote:   cacheFriendship.TieStrength,
						},
					)
				} else {
					pointsAllow -= cacheFriendship.TieStrength
					cacheConflict.Reasoning = append(
						cacheConflict.Reasoning,
						domain.Reason{
							UserID: partyUser.ID,
							Vote:   -cacheFriendship.TieStrength,
						},
					)
				}
			}

			if pointsAllow > 0 {
				cacheConflict.Result = "allow"
			} else if pointsAllow < 0 {
				cacheConflict.Result = "block"
			} else {
				cacheConflict.Result = "indeterminate"
			}
			cachePhoto.Conflicts = append(cachePhoto.Conflicts, cacheConflict)
		}
	}

	c.photoRedis.Save(cachePhoto)

	for _, u := range dbPhoto.TaggedUsers {
		c.userRedis.Publish(u.ID, "photo", domain.WebPhotoFromCachePhoto(cachePhoto))
	}

	elapsed := time.Since(start)
	c.logger.Printf("Processed photo %s conflicts in %s", dbPhoto.ID, elapsed)
}

func shouldAllow(
	taggedUserAllowedUsers map[string][]string,
	taggedUserBlockedUsers map[string][]string,
	partyUserID string,
	targetUserID string,
) bool {
	for _, allowedUserID := range taggedUserAllowedUsers[partyUserID] {
		if targetUserID == allowedUserID {
			return true
		}
	}

	for _, blockedUserID := range taggedUserBlockedUsers[partyUserID] {
		if targetUserID == blockedUserID {
			return false
		}
	}

	log.Printf("warning: Could not find preference for %s with %s", partyUserID, targetUserID)
	return true // default true as abstained users do not care
}

func isIn(needle domain.DBCategory, haystack []domain.DBCategory) bool {
	for _, cat := range haystack {
		if needle.Name == cat.Name {
			return true
		}
	}
	return false
}

func isInUsers(needle string, haystack []domain.DBUser) bool {
	for _, u := range haystack {
		if u.ID == needle {
			return true
		}
	}

	return false
}
