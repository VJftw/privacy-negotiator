package survey

import (
	"log"

	"encoding/json"
	"time"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/utils"
	"github.com/streadway/amqp"
)

// Consumer - Consumer for persisting Surveys.
type Consumer struct {
	queue    amqp.Queue
	channel  *amqp.Channel
	logger   *log.Logger
	surveyDB *DBManager
}

// NewConsumer - Returns a new consumer.
func NewConsumer(
	queueLogger *log.Logger,
	ch *amqp.Channel,
	surveyDBManager *DBManager,
) *Consumer {
	queue, err := ch.QueueDeclare(
		"persist-survey", // name
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")
	return &Consumer{
		logger:   queueLogger,
		channel:  ch,
		queue:    queue,
		surveyDB: surveyDBManager,
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

func (c *Consumer) process(d amqp.Delivery) {
	start := time.Now()

	dbSurvey := &domain.DBSurvey{}
	json.Unmarshal(d.Body, dbSurvey)

	c.logger.Printf("Started processing survey for %s", dbSurvey.UserID)

	c.surveyDB.Save(dbSurvey)

	c.logger.Printf("Processed survey for %s in %s", dbSurvey.UserID, time.Since(start))
	d.Ack(false)
}
