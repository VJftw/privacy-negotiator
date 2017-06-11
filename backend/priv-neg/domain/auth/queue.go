package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/utils"
	"github.com/streadway/amqp"
)

// AuthQueue - Queue for publishing jobs to get long-lived facebook tokens.
type AuthQueue struct {
	queue       amqp.Queue
	channel     *amqp.Channel
	userManager user.Managable
	logger      *log.Logger
}

// NewAuthQueue - Returns an implementation of the queue.
func NewAuthQueue(queueLogger *log.Logger, userManager user.Managable) *AuthQueue {
	return &AuthQueue{
		logger:      queueLogger,
		userManager: userManager,
	}
}

// Setup - Declares the Queue.
func (q *AuthQueue) Setup(ch *amqp.Channel) {
	queue, err := ch.QueueDeclare(
		"get-facebook-long-lived-token", // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")
	q.queue = queue
	q.channel = ch
}

// Publish - Does nothing in the Worker.
func (q *AuthQueue) Publish(i utils.Queueable) {
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

// Consume - Processes items from the Queue.
func (q *AuthQueue) Consume() {
	msgs, err := q.channel.Consume(
		q.queue.Name, // queue
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
			q.process(d)
		}
	}()

	q.logger.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (q *AuthQueue) process(d amqp.Delivery) {
	start := time.Now()

	user := q.userManager.New()
	json.Unmarshal(d.Body, user)

	q.logger.Printf("Started processing for %s", user.FacebookUserID)

	respLongLived := getLongLivedToken(user)

	user.LongLivedToken = respLongLived.AccessToken
	duration, _ := time.ParseDuration(fmt.Sprintf("%ds", respLongLived.Expires))
	user.TokenExpires = time.Now().Add(duration)

	q.userManager.Save(user)

	elapsed := time.Since(start)
	q.logger.Printf("Processed AuthQueue for %s in %s", user.FacebookUserID, elapsed)
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
