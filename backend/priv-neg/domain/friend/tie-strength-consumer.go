package friend

import (
	"log"
	"time"

	"encoding/json"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/utils"
	"github.com/streadway/amqp"
)

// TieStrengthConsumer - Consumer for getting Community Detection.
type TieStrengthConsumer struct {
	queue       amqp.Queue
	channel     *amqp.Channel
	logger      *log.Logger
	userRedis   *user.RedisManager
	friendRedis *RedisManager
}

// NewTieStrengthConsumer - Returns a new consumer.
func NewTieStrengthConsumer(
	queueLogger *log.Logger,
	ch *amqp.Channel,
	userRedisManager *user.RedisManager,
	friendRedisManager *RedisManager,
) *TieStrengthConsumer {
	queue, err := ch.QueueDeclare(
		"tie-strength", // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")
	return &TieStrengthConsumer{
		logger:      queueLogger,
		channel:     ch,
		queue:       queue,
		userRedis:   userRedisManager,
		friendRedis: friendRedisManager,
	}
}

// Consume - Processes items from the Queue.
func (c *TieStrengthConsumer) Consume() {
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

func (c *TieStrengthConsumer) process(d amqp.Delivery) {
	start := time.Now()

	queueFriendship := domain.QueueFriendship{}
	json.Unmarshal(d.Body, &queueFriendship)
	c.logger.Printf("Starting processing for %s:%s", queueFriendship.From, queueFriendship.To)

	// Go through both user's profiles (in cache) and compare similarities. Similarity is bi-directional!
	aCacheProfile, err := c.userRedis.GetProfileByID(queueFriendship.From)
	if err != nil {
		c.logger.Printf("User profile not present in Cache %s", queueFriendship.From)
		return
	}
	bCacheProfile, err := c.userRedis.GetProfileByID(queueFriendship.To)
	if err != nil {
		c.logger.Printf("User profile not present in Cache %s", queueFriendship.To)
		return
	}

	points := 0

	// Gender
	if aCacheProfile.Gender == bCacheProfile.Gender {
		points++
	}

	// AgeRange
	if aCacheProfile.AgeRange == bCacheProfile.AgeRange {
		points++
	}

	// Hometown
	if aCacheProfile.Hometown == bCacheProfile.Hometown {
		points++
	}

	// Location
	if aCacheProfile.Location == bCacheProfile.Location {
		points++
	}

	// Educations
	for _, v := range aCacheProfile.Education {
		if utils.IsIn(v, bCacheProfile.Education) {
			points++
		}
	}
	// Favorite Teams
	for _, v := range aCacheProfile.FavouriteTeams {
		if utils.IsIn(v, bCacheProfile.FavouriteTeams) {
			points++
		}
	}

	// Inspirational People
	for _, v := range aCacheProfile.InspirationalPeople {
		if utils.IsIn(v, bCacheProfile.InspirationalPeople) {
			points++
		}
	}

	// Languages
	for _, v := range aCacheProfile.Languages {
		if utils.IsIn(v, bCacheProfile.Languages) {
			points++
		}
	}

	// Music
	for _, v := range aCacheProfile.Music {
		if utils.IsIn(v, bCacheProfile.Music) {
			points++
		}
	}

	// Movies
	for _, v := range aCacheProfile.Movies {
		if utils.IsIn(v, bCacheProfile.Movies) {
			points++
		}
	}

	// Likes
	for _, v := range aCacheProfile.Likes {
		if utils.IsIn(v, bCacheProfile.Likes) {
			points++
		}
	}

	// Groups
	for _, v := range aCacheProfile.Groups {
		if utils.IsIn(v, bCacheProfile.Groups) {
			points++
		}
	}

	// Events
	for _, v := range aCacheProfile.Events {
		if utils.IsIn(v, bCacheProfile.Events) {
			points++
		}
	}

	// TODO: Family

	aCacheUser, _ := c.userRedis.FindByID(queueFriendship.From)
	aCacheFriendship, _ := c.friendRedis.FindByIDAndUser(queueFriendship.To, aCacheUser)
	aCacheFriendship.TieStrength = points

	bCacheUser, _ := c.userRedis.FindByID(queueFriendship.To)
	bCacheFriendship, err := c.friendRedis.FindByIDAndUser(queueFriendship.From, bCacheUser)
	if err != nil {
		bCacheFriendship = &domain.CacheFriendship{
			ID:          aCacheUser.ID,
			TieStrength: points,
		}
	}

	c.friendRedis.Save(aCacheUser, aCacheFriendship)
	c.friendRedis.Save(bCacheUser, bCacheFriendship)

	c.userRedis.Publish(aCacheUser.ID, "clique", aCacheFriendship)
	c.userRedis.Publish(bCacheUser.ID, "clique", bCacheFriendship)
	c.logger.Printf("Processed tie-strength %s:%s in %s", queueFriendship.From, queueFriendship.To, time.Since(start))
}
