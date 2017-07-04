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
	queue         amqp.Queue
	channel       *amqp.Channel
	logger        *log.Logger
	categoryDB    *DBManager
	categoryRedis *RedisManager
}

// NewConsumer - Returns a new consumer.
func NewConsumer(
	queueLogger *log.Logger,
	ch *amqp.Channel,
	categoryDBManager *DBManager,
	categoryRedisManager *RedisManager,
) *Consumer {
	queue, err := ch.QueueDeclare(
		"persist-category", // name
		true,               // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")
	return &Consumer{
		logger:        queueLogger,
		channel:       ch,
		queue:         queue,
		categoryDB:    categoryDBManager,
		categoryRedis: categoryRedisManager,
	}
}

func (c *Consumer) loadCategories() {
	categories := []string{
		"Family",
		"Party",
		"NSFW (Not Safe For Work)",
		"Holiday",
		"Friends",
	}

	for _, nameCategory := range categories {
		c.categoryRedis.Save(nameCategory)
		_, err := c.categoryDB.FindByName(nameCategory)
		if err != nil {
			c.categoryDB.Save(&domain.DBCategory{Name: nameCategory, UserID: "none"})
		}
	}
}

// Consume - Processes items from the Queue.
func (c *Consumer) Consume() {
	c.loadCategories()
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

	queueCategory := &domain.QueueCategory{}
	json.Unmarshal(d.Body, queueCategory)

	c.logger.Printf("Started processing for category %s", queueCategory.Name)

	dbCategory := domain.DBCategoryFromQueueCategory(queueCategory)

	c.categoryDB.Save(dbCategory)

	elapsed := time.Since(start)
	c.logger.Printf("Processed Category %s in %s", dbCategory.Name, elapsed)
}
