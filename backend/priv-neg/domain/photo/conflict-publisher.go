package photo

import (
	"encoding/json"
	"log"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/persisters"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/utils"
	"github.com/streadway/amqp"
)

// ConflictPublisher - Publishes messages to the queue for getting tagged users in a photo.
type ConflictPublisher struct {
	queue   amqp.Queue
	channel *amqp.Channel
	logger  *log.Logger
}

// NewConflictPublisher - Returns a new ConflictPublisher.
func NewConflictPublisher(
	queueLogger *log.Logger,
	ch *amqp.Channel,
) *ConflictPublisher {
	queue, err := ch.QueueDeclare(
		"conflict-detection-and-resolution", // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")
	return &ConflictPublisher{
		logger:  queueLogger,
		channel: ch,
		queue:   queue,
	}
}

// Publish - Publishes a given message onto the Queue.
func (q *ConflictPublisher) Publish(i persisters.Queueable) {
	b, err := json.Marshal(i)
	if err != nil {
		log.Println(err)
		return
	}

	err = q.channel.Publish(
		"",           // exchange
		q.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         b,
		})
	utils.FailOnError(err, "Failed to publish a message")
}

// GetMessageTotal - returns the total amount of messages in the queue.
func (q *ConflictPublisher) GetMessageTotal() int {
	return q.queue.Messages
}
