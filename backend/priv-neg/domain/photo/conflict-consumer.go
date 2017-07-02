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
	queue      amqp.Queue
	channel    *amqp.Channel
	logger     *log.Logger
	friendDB   *friend.DBManager
	userDB     *user.DBManager
	photoDB    *DBManager
	photoRedis *RedisManager
	userRedis  *user.RedisManager
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
		logger:     queueLogger,
		channel:    ch,
		queue:      queue,
		friendDB:   friendDBManager,
		userDB:     userDBManager,
		photoDB:    photoDBManager,
		photoRedis: photoRedisManager,
		userRedis:  userRedisManager,
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

	c.logger.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (c *ConflictConsumer) process(d amqp.Delivery) {
	start := time.Now()

	dbPhoto := domain.DBPhoto{}
	json.Unmarshal(d.Body, &dbPhoto)

	taggedUserAllowedUsers := map[string][]string{} // taggedUserID: cliqueUserID
	taggedUserBlockedUsers := map[string][]string{} // taggedUserID: cliqueUserID

	c.logger.Printf("DEBUG: TaggedUsers: %d, Categories: %d", len(dbPhoto.TaggedUsers), len(dbPhoto.Categories))
	taggedUserIDs := []string{}
	for _, taggedUser := range dbPhoto.TaggedUsers {
		taggedUserIDs = append(taggedUserIDs, taggedUser.ID)
	}

	for _, taggedUser := range dbPhoto.TaggedUsers {
		userInitialAllowed := []string{taggedUser.ID}
		userInitialBlocked := []string{}
		taggedUserCliques, _ := c.friendDB.GetUserCliquesByUser(taggedUser)
		for _, taggedUserClique := range taggedUserCliques {
			c.logger.Printf("DEBUG: UserCliqueCategories: %v", taggedUserClique.Categories)
			hasCategory := false
			for _, cliqueCat := range taggedUserClique.Categories {
				c.logger.Printf("DEBUG: Allowing UserClique: %s", taggedUserClique.CliqueID)
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

	c.logger.Printf("DEBUG: Allowed Users %v", taggedUserAllowedUsers)
	c.logger.Printf("DEBUG: Blocked Users %v", taggedUserBlockedUsers)

	// Find Conflicts
	dbConflict := domain.NewDBConflict()
	dbConflict.Photo = dbPhoto
	dbConflict.PhotoID = dbPhoto.ID

	cachePhoto, _ := c.photoRedis.FindByID(dbPhoto.ID)
	cachePhoto.AllowedUserIDs = []string{}
	cachePhoto.BlockedUserIDs = []string{}
	// Update AllowedUsers and BlockedUsers for CachedPhoto as well as finding conflicts.
	for taggedUserIDAllowed, allowedUserIDs := range taggedUserAllowedUsers {
		for _, allowedUserID := range allowedUserIDs {
			conflict := false
			for taggedUserIDBlocked, blockedUserIDs := range taggedUserBlockedUsers {
				for _, blockedUserID := range blockedUserIDs {
					if allowedUserID == blockedUserID {
						// conflict
						c.logger.Printf("Found conflict between %s and %s with %s", taggedUserIDAllowed, taggedUserIDBlocked, allowedUserID)
						allowedUser, _ := c.userDB.FindByID(allowedUserID)
						dbConflict.Targets = append(dbConflict.Targets, *allowedUser)
						taggedUserAllowed, _ := c.userDB.FindByID(taggedUserIDAllowed)
						taggedUserBlocked, _ := c.userDB.FindByID(taggedUserIDBlocked)
						dbConflict.Parties = append(dbConflict.Parties, *taggedUserBlocked, *taggedUserAllowed)
						conflict = true
					} else if !utils.IsIn(blockedUserID, cachePhoto.BlockedUserIDs) {
						cachePhoto.BlockedUserIDs = append(cachePhoto.BlockedUserIDs, blockedUserID)
					}
				}
			}
			if !conflict && !utils.IsIn(allowedUserID, cachePhoto.AllowedUserIDs) {
				cachePhoto.AllowedUserIDs = append(cachePhoto.AllowedUserIDs, allowedUserID)
			}
		}
	}

	cachePhoto.Conflict = domain.CacheConflict{}
	if len(dbConflict.Targets) > 0 {
		// We have conflicts. Save to Cache and DB and suggest resolution
		cacheConflict := domain.CacheConflictFromDBConflict(dbConflict)
		cachePhoto.Conflict = cacheConflict

		err := c.photoDB.SaveConflict(&dbConflict)
		if err != nil {
			c.logger.Printf("ERROR: %v", err)
		}

		// For every target user
		// Find which party user has the highest tie-strength to target user.
		// Suggest policy of that party user

	}

	c.photoRedis.Save(cachePhoto)

	for _, u := range dbPhoto.TaggedUsers {
		c.userRedis.Publish(u.ID, "photo", domain.WebPhotoFromCachePhoto(cachePhoto))
	}

	elapsed := time.Since(start)
	c.logger.Printf("Processed photo %s conflicts in %s", dbPhoto.ID, elapsed)
}

func isIn(needle domain.DBCategory, haystack []domain.DBCategory) bool {
	for _, cat := range haystack {
		if needle.Name == cat.Name {
			return true
		}
	}
	return false
}
