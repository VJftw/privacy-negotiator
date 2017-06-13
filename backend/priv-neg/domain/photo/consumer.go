package photo

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/utils"
	"github.com/streadway/amqp"
)

// Consumer - Consumer for getting TaggedUsers for a Photo.
type Consumer struct {
	queue      amqp.Queue
	channel    *amqp.Channel
	logger     *log.Logger
	userDB     *user.DBManager
	photoRedis *RedisManager
	userRedis  *user.RedisManager
}

// NewConsumer - Returns a new consumer.
func NewConsumer(
	queueLogger *log.Logger,
	ch *amqp.Channel,
	userDBManager *user.DBManager,
	userRedisManager *user.RedisManager,
	photoRedisManager *RedisManager,
) *Consumer {
	queue, err := ch.QueueDeclare(
		"photo-tags", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")
	return &Consumer{
		logger:     queueLogger,
		channel:    ch,
		queue:      queue,
		userDB:     userDBManager,
		userRedis:  userRedisManager,
		photoRedis: photoRedisManager,
	}
}

// Consume - Processes items from the Queue.
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

	cachePhoto := &domain.CachePhoto{}
	json.Unmarshal(d.Body, cachePhoto)

	c.logger.Printf("Started processing for %s", cachePhoto.ID)

	dbUser, _ := c.userDB.FindByID(cachePhoto.Uploader)

	updatePhotoFromGraphAPI(cachePhoto, dbUser)

	c.photoRedis.Save(cachePhoto)

	webPhoto := domain.WebPhotoFromCachePhoto(cachePhoto)
	for _, user := range cachePhoto.TaggedUsers {
		c.userRedis.Publish(user, "photo", webPhoto)
	}

	elapsed := time.Since(start)
	c.logger.Printf("Processed SyncPhoto for %s in %s", cachePhoto.ID, elapsed)
}

func updatePhotoFromGraphAPI(p *domain.CachePhoto, u *domain.DBUser) {

	res, _ := http.Get(fmt.Sprintf(
		"https://graph.facebook.com/v2.9/%s?access_token=%s&fields=from,tags{id}",
		p.ID,
		u.LongLivedToken,
	))

	photoResponse := &fbResponsePhoto{}
	err := json.NewDecoder(res.Body).Decode(photoResponse)
	defer res.Body.Close()
	if err != nil {
		log.Printf("Error: %s", err)
	}
	p.Uploader = photoResponse.From.ID

	for _, taggedUser := range photoResponse.Tags.Data {
		p.TaggedUsers = append(p.TaggedUsers, taggedUser.ID)
	}

	p.Pending = false

}

type fbResponseUser struct {
	ID string `json:"id"`
}

type fbResponsePhoto struct {
	From fbResponseUser  `json:"from"`
	Tags fbResponsePager `json:"tags"`
}

type fbResponsePager struct {
	Data []fbResponseUser `json:"data"`
}
