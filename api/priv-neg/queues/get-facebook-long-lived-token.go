package queues

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type GetFacebookLongLivedToken struct {
	queue   amqp.Queue
	channel *amqp.Channel
}

func NewGetFacebookLongLivedToken() *GetFacebookLongLivedToken {
	return &GetFacebookLongLivedToken{}
}

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

func (q *GetFacebookLongLivedToken) Consume() {

}
