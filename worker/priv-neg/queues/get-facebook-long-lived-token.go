package queues

import (
	"encoding/json"
	"log"

	"github.com/VJftw/privacy-negotiator/worker/priv-neg/domain/user"
	"github.com/streadway/amqp"
)

type GetFacebookLongLivedToken struct {
	queue       amqp.Queue
	channel     *amqp.Channel
	UserManager user.Manager `inject:"user.manager"`
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
	q.Consume()
}

func (q *GetFacebookLongLivedToken) Publish(i Queueable) {
}

func (q *GetFacebookLongLivedToken) Consume() {
	msgs, err := q.channel.Consume(
		q.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			q.process(d)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (q *GetFacebookLongLivedToken) process(d amqp.Delivery) {
	log.Printf("Received a message: %s", d.Body)
	user := q.UserManager.New()

	json.Unmarshal(d.Body, user)

	log.Printf("Created user: %s", user)

}
