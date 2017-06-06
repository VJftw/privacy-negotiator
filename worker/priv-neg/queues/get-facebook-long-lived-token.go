package queues

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
	start := time.Now()

	user := q.UserManager.New()
	json.Unmarshal(d.Body, user)

	log.Printf("Started processing for %s\n", user.FacebookUserID)

	respLongLived := GetLongLivedToken(user)

	user.LongLivedToken = respLongLived.AccessToken
	duration, _ := time.ParseDuration(fmt.Sprintf("%ds", respLongLived.Expires))
	user.TokenExpires = time.Now().Add(duration)

	q.UserManager.Save(user)

	elapsed := time.Since(start)
	log.Printf("Processed GetFacebookLongLivedToken for %s in %s\n", user.FacebookUserID, elapsed)
}

func GetLongLivedToken(fbUser *user.FacebookUser) *facebookResponseLongLived {
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
