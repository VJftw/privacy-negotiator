package queues

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

// GetFacebookLongLivedToken - Queue for publishing jobs to get long-lived facebook tokens.
type GetFacebookLongLivedToken struct {
	queue   amqp.Queue
	channel *amqp.Channel
}

// NewGetFacebookLongLivedToken - Returns an implementation of the queue.
func NewGetFacebookLongLivedToken() *GetFacebookLongLivedToken {
	return &GetFacebookLongLivedToken{}
}

// Setup - Declares the Queue.
func (q *GetFacebookLongLivedToken) Setup(ch *amqp.Channel) {
	queue, err := ch.QueueDeclare(
		"get-facebook-long-lived-token", // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")
	q.queue = queue
	q.channel = ch
}

// Publish - Adds an item to the queue.
func (q *GetFacebookLongLivedToken) Publish(i Queueable) {
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
	failOnError(err, "Failed to publish a message")
}

// Consume - Does nothing in the API.
func (q *GetFacebookLongLivedToken) Consume() {}
