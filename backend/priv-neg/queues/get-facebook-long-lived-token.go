package queues

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
	"github.com/streadway/amqp"
)

// GetFacebookLongLivedToken - Queue for publishing jobs to get long-lived facebook tokens.
type GetFacebookLongLivedToken struct {
	queue       amqp.Queue
	channel     *amqp.Channel
	UserManager user.Managable `inject:"user.manager"`
	Logger      *log.Logger    `inject:"logger.queue"`
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

// Publish - Does nothing in the Worker.
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

// Consume - Processes items from the Queue.
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

	q.Logger.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (q *GetFacebookLongLivedToken) process(d amqp.Delivery) {
	start := time.Now()

	user := q.UserManager.New()
	json.Unmarshal(d.Body, user)

	q.Logger.Printf("Started processing for %s\n", user.FacebookUserID)

	respLongLived := getLongLivedToken(user)

	user.LongLivedToken = respLongLived.AccessToken
	duration, _ := time.ParseDuration(fmt.Sprintf("%ds", respLongLived.Expires))
	user.TokenExpires = time.Now().Add(duration)

	q.UserManager.Save(user)

	elapsed := time.Since(start)
	q.Logger.Printf("Processed GetFacebookLongLivedToken for %s in %s\n", user.FacebookUserID, elapsed)
}

func getLongLivedToken(fbUser *user.FacebookUser) *facebookResponseLongLived {
	clientID := os.Getenv("FACEBOOK_APP_ID")
	clientSecret := os.Getenv("FACEBOOK_APP_SECRET")
	res, _ := http.Get(fmt.Sprintf("https://graph.facebook.com/v2.9/oauth/access_token?grant_type=fb_exchange_token&client_id=%s&client_secret=%s&fb_exchange_token=%s",
		clientID,
		clientSecret,
		fbUser.ShortLivedToken,
	))

	respLongLived := &facebookResponseLongLived{}
	_ = json.NewDecoder(res.Body).Decode(respLongLived)

	return respLongLived

}

type facebookResponseLongLived struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Expires     uint   `json:"expires_in"`
}
