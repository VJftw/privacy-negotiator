package photo

import (
	"encoding/json"
	"log"
	"time"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/utils"
	"github.com/streadway/amqp"
)

// PersistConsumer - Consumer for getting Community Detection.
type PersistConsumer struct {
	queue             amqp.Queue
	channel           *amqp.Channel
	logger            *log.Logger
	photoDB           *DBManager
	userDB            *user.DBManager
	conflictPublisher *ConflictPublisher
}

// NewPersistConsumer - Returns a new consumer.
func NewPersistConsumer(
	queueLogger *log.Logger,
	ch *amqp.Channel,
	photoDBManager *DBManager,
	userDBManager *user.DBManager,
	conflictPublisher *ConflictPublisher,
) *PersistConsumer {
	queue, err := ch.QueueDeclare(
		"persist-photo", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")
	return &PersistConsumer{
		logger:            queueLogger,
		channel:           ch,
		queue:             queue,
		photoDB:           photoDBManager,
		userDB:            userDBManager,
		conflictPublisher: conflictPublisher,
	}
}

// Consume - Processes items from the Queue.
func (c *PersistConsumer) Consume() {
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

func (c *PersistConsumer) process(d amqp.Delivery) {
	start := time.Now()

	dbPhoto := domain.DBPhoto{}
	json.Unmarshal(d.Body, &dbPhoto)

	dbUsers := []domain.DBUser{}
	for _, user := range dbPhoto.TaggedUsers {
		dbUser, err := c.userDB.FindByID(user.ID)
		if err == nil {
			dbUsers = append(dbUsers, *dbUser)
		}
	}
	dbPhoto.TaggedUsers = dbUsers
	dbUploader, _ := c.userDB.FindByID(dbPhoto.Uploader)
	dbPhoto.Uploader = dbUploader.ID

	c.photoDB.Save(&dbPhoto)

	c.conflictPublisher.Publish(dbPhoto)

	elapsed := time.Since(start)
	c.logger.Printf("Processed photo %s in %s", dbPhoto.ID, elapsed)
}
