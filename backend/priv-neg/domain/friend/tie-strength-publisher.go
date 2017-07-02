package friend

import (
	"encoding/json"
	"log"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/persisters"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/utils"
	"github.com/streadway/amqp"
)

// TieStrengthPublisher - Publishes messages to the queue for performing community detection.
type TieStrengthPublisher struct {
	queue   amqp.Queue
	channel *amqp.Channel
	logger  *log.Logger
}

// NewTieStrengthPublisher - Returns a new TieStrengthPublisher.
func NewTieStrengthPublisher(
	queueLogger *log.Logger,
	ch *amqp.Channel,
) *TieStrengthPublisher {
	queue, err := ch.QueueDeclare(
		"tie-strength", // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")
	return &TieStrengthPublisher{
		logger:  queueLogger,
		channel: ch,
		queue:   queue,
	}
}

// Publish - Publishes a given message onto the Queue.
func (q *TieStrengthPublisher) Publish(i persisters.Queueable) {
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
func (q *TieStrengthPublisher) GetMessageTotal() int {
	return q.queue.Messages
}
