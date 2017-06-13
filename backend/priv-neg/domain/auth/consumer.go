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

// Consumer - Consumer for authentication. Gets a long-lived Facebook API token.
type Consumer struct {
	queue   amqp.Queue
	channel *amqp.Channel
	logger  *log.Logger
	userDB  *user.DBManager
}

// NewConsumer - Returns a new consumer for authentication.
func NewConsumer(
	queueLogger *log.Logger,
	ch *amqp.Channel,
	userDBManager *user.DBManager,
) *Consumer {
	queue, err := ch.QueueDeclare(
		"auth-long-token", // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")
	return &Consumer{
		logger:  queueLogger,
		channel: ch,
		queue:   queue,
		userDB:  userDBManager,
	}
}

// Consume - Starts running the consumer.
func (c *Consumer) Consume() {
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

	authUser := &user.AuthUser{}
	json.Unmarshal(d.Body, authUser)

	c.logger.Printf("Started processing for %s", authUser.ID)

	respLongLived := getLongLivedToken(authUser)

	dbUser := user.DBUserFromAuthUser(authUser)
	dbUser.LongLivedToken = respLongLived.AccessToken
	duration, _ := time.ParseDuration(fmt.Sprintf("%ds", respLongLived.Expires))
	dbUser.TokenExpires = time.Now().Add(duration)

	c.userDB.Save(dbUser)

	elapsed := time.Since(start)
	c.logger.Printf("Processed AuthQueue for %s in %s", dbUser.ID, elapsed)
}

func getLongLivedToken(u *user.AuthUser) *facebookResponseLongLived {
	clientID := os.Getenv("FACEBOOK_APP_ID")
	clientSecret := os.Getenv("FACEBOOK_APP_SECRET")
	res, _ := http.Get(fmt.Sprintf("https://graph.facebook.com/v2.9/oauth/access_token?grant_type=fb_exchange_token&client_id=%s&client_secret=%s&fb_exchange_token=%s",
		clientID,
		clientSecret,
		u.ShortLivedToken,
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
