package conflict

import (
	"encoding/json"
	"log"
	"time"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/friend"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/utils"
	"github.com/streadway/amqp"
)

// Consumer - Consumer for conflict detection and suggested resolution.
type Consumer struct {
	queue    amqp.Queue
	channel  *amqp.Channel
	logger   *log.Logger
	friendDB *friend.DBManager
}

// NewConsumer - Returns a new consumer.
func NewConsumer(
	queueLogger *log.Logger,
	ch *amqp.Channel,
	friendDBManager *friend.DBManager,
) *Consumer {
	queue, err := ch.QueueDeclare(
		"conflict-detection-and-resolution", // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")
	return &Consumer{
		logger:   queueLogger,
		channel:  ch,
		queue:    queue,
		friendDB: friendDBManager,
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

	c.logger.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (c *Consumer) process(d amqp.Delivery) {
	start := time.Now()

	dbPhoto := domain.DBPhoto{}
	json.Unmarshal(d.Body, &dbPhoto)

	taggedUserAllowedUsers := map[string]map[string]bool{} // taggedUserID: cliqueUserID
	taggedUserBlockedUsers := map[string]map[string]bool{} // taggedUserID: cliqueUserID

	c.logger.Printf("DEBUG: TaggedUsers: %d, Categories: %d", len(dbPhoto.TaggedUsers), len(dbPhoto.Categories))

	for _, taggedUser := range dbPhoto.TaggedUsers { // for each tagged user
		taggedUserAllowedUsers[taggedUser.ID] = map[string]bool{}
		taggedUserBlockedUsers[taggedUser.ID] = map[string]bool{}

		dbUserCliques, _ := c.friendDB.GetUserCliquesByUser(taggedUser)
		for _, userClique := range dbUserCliques { // for each of the user's cliques
			// if the clique has one of the categories of the photo, then add every user to allowed. If not, then add to blocked.
			hasCategory := false
			c.logger.Printf("DEBUG: User Clique %s Categories: %v", userClique.Name, userClique.Categories)
			for _, cliqueCat := range userClique.Categories {
				c.logger.Printf("CliqueCat: %v", cliqueCat)
				if isIn(cliqueCat, dbPhoto.Categories) {
					hasCategory = true
					break
				}
			}
			clique, _ := c.friendDB.FindCliqueByID(userClique.CliqueID)
			if hasCategory {
				for _, userID := range clique.GetUserIDs() {
					taggedUserAllowedUsers[taggedUser.ID][userID] = true
				}
			} else {
				for _, userID := range clique.GetUserIDs() {
					taggedUserBlockedUsers[taggedUser.ID][userID] = true
				}
			}
		}
	}

	c.logger.Printf("DEBUG: Allowed Users %v", taggedUserAllowedUsers)
	c.logger.Printf("DEBUG: Blocked Users %v", taggedUserBlockedUsers)

	allowedUsers := []string{}
	blockedUsers := []string{}

	for taggedUserID, allowedUserIDs := range taggedUserAllowedUsers {
		allowedUsers = append(allowedUsers, taggedUserID)
		for allowedUserID := range allowedUserIDs {
			allowedUsers = append(allowedUsers, allowedUserID)
		}
		if _, ok := taggedUserBlockedUsers[taggedUserID]; ok {
			for blockedUserID := range taggedUserBlockedUsers[taggedUserID] {
				if _, ok := taggedUserAllowedUsers[taggedUserID][blockedUserID]; !ok {
					blockedUsers = append(blockedUsers, blockedUserID)
				}
			}
		}
	}

	c.logger.Printf("DEBUG: Allowed Users %v", allowedUsers)
	c.logger.Printf("DEBUG: Blocked Users %v", blockedUsers)

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
