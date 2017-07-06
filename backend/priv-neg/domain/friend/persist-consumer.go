package friend

import (
	"log"
	"time"

	"encoding/json"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/utils"
	"github.com/streadway/amqp"
)

// PersistConsumer - Consumer for getting Community Detection.
type PersistConsumer struct {
	queue    amqp.Queue
	channel  *amqp.Channel
	logger   *log.Logger
	cliqueDB *DBManager
}

// NewPersistConsumer - Returns a new consumer.
func NewPersistConsumer(
	queueLogger *log.Logger,
	ch *amqp.Channel,
	cliqueDBManager *DBManager,
) *PersistConsumer {
	queue, err := ch.QueueDeclare(
		"persist-user-clique", // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")
	return &PersistConsumer{
		logger:   queueLogger,
		channel:  ch,
		queue:    queue,
		cliqueDB: cliqueDBManager,
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

	c.logger.Printf("Worker: %s waiting for messages. To exit press CTRL+C", c.queue.Name)
	<-forever
}

func (c *PersistConsumer) process(d amqp.Delivery) {
	start := time.Now()

	dbUserClique := domain.DBUserClique{}
	json.Unmarshal(d.Body, &dbUserClique)
	c.logger.Printf("debug: UserClique: %v", dbUserClique)

	c.cliqueDB.SaveUserClique(&dbUserClique)

	elapsed := time.Since(start)
	c.logger.Printf("Processed user clique %s:%s in %s", dbUserClique.UserID, dbUserClique.CliqueID, elapsed)
}
