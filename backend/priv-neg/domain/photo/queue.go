package photo

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/utils"
	"github.com/streadway/amqp"
)

// SyncQueue - Queue for syncing Photo state.
type SyncQueue struct {
	logger       *log.Logger
	photoManager Managable
	userManager  user.Managable
	queue        amqp.Queue
	channel      *amqp.Channel
}

// NewSyncQueue - Returns a new queue for syncing photos.
func NewSyncQueue(
	logger *log.Logger,
	photoManager Managable,
	userManager user.Managable,
) *SyncQueue {
	return &SyncQueue{
		logger:       logger,
		photoManager: photoManager,
		userManager:  userManager,
	}
}

// Setup - Declares the Queue.
func (q *SyncQueue) Setup(ch *amqp.Channel) {
	queue, err := ch.QueueDeclare(
		"photo-sync", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")
	q.queue = queue
	q.channel = ch
}

// Publish - Publishes a message to the queue.
func (q *SyncQueue) Publish(i utils.Queueable) {
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
func (q *SyncQueue) Consume() {
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

func (q *SyncQueue) process(d amqp.Delivery) {
	start := time.Now()

	photo := &FacebookPhoto{}
	json.Unmarshal(d.Body, photo)

	q.logger.Printf("Started processing for %s", photo.FacebookPhotoID)

	user, _ := q.userManager.FindByID(photo.Uploader)

	updatePhotoFromGraphAPI(photo, user)

	q.photoManager.Save(photo)

	elapsed := time.Since(start)
	q.logger.Printf("Processed SyncPhoto for %s in %s", photo.FacebookPhotoID, elapsed)
}

func updatePhotoFromGraphAPI(p *FacebookPhoto, authUser *user.FacebookUser) {

	res, _ := http.Get(fmt.Sprintf(
		"https://graph.facebook.com/v2.9/%s?access_token=%s&fields=from,tags{id}",
		p.FacebookPhotoID,
		authUser.LongLivedToken,
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
