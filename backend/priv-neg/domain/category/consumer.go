package category

import (
	"log"

	"encoding/json"
	"time"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/utils"
	"github.com/streadway/amqp"
)

// Consumer - Consumer for persisting Categories.
type Consumer struct {
	queue      amqp.Queue
	channel    *amqp.Channel
	logger     *log.Logger
	categoryDB *DBManager
}

// NewConsumer - Returns a new consumer.
func NewConsumer(
	queueLogger *log.Logger,
	ch *amqp.Channel,
	categoryDBManager *DBManager,
) *Consumer {
	queue, err := ch.QueueDeclare(
		"category-persist", // name
		true,               // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")
	return &Consumer{
		logger:     queueLogger,
		channel:    ch,
		queue:      queue,
		categoryDB: categoryDBManager,
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

	dbCategory := &domain.DBCategory{}
	json.Unmarshal(d.Body, dbCategory)

	c.logger.Printf("Started processing for category %s:%s", dbCategory.UserID, dbCategory.Name)

	c.categoryDB.Save(dbCategory)

	elapsed := time.Since(start)
	c.logger.Printf("Processed Category %s:%s in %s", dbCategory.UserID, dbCategory.Name, elapsed)
}
