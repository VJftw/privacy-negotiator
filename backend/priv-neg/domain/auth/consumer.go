package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"strconv"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/friend"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/utils"
	"github.com/streadway/amqp"
)

// Consumer - Consumer for authentication. Gets a long-lived Facebook API token.
type Consumer struct {
	queue           amqp.Queue
	channel         *amqp.Channel
	logger          *log.Logger
	userDB          *user.DBManager
	friendPublisher *friend.Publisher
	userRedis       *user.RedisManager
}

// NewConsumer - Returns a new consumer for authentication.
func NewConsumer(
	queueLogger *log.Logger,
	ch *amqp.Channel,
	userDBManager *user.DBManager,
	friendPublisher *friend.Publisher,
	userRedisManager *user.RedisManager,
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
		logger:          queueLogger,
		channel:         ch,
		queue:           queue,
		userDB:          userDBManager,
		friendPublisher: friendPublisher,
		userRedis:       userRedisManager,
	}
}

// Consume - Starts running the consumer.
func (c *Consumer) Consume() {
	msgs, err := c.channel.Consume(
		c.queue.Name, // queue
		"",           // consumer
		false,        // auto-ack
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

	c.logger.Printf("Worker: %s waiting for messages. To exit press CTRL+C", c.queue.Name)
	<-forever
}

func (c *Consumer) process(d amqp.Delivery) {
	start := time.Now()

	authUser := &domain.AuthUser{}
	json.Unmarshal(d.Body, authUser)

	c.logger.Printf("Started processing for %s", authUser.ID)

	respLongLived := getLongLivedToken(authUser)

	dbUser := domain.DBUserFromAuthUser(authUser)
	dbUser.LongLivedToken = respLongLived.AccessToken
	duration, _ := time.ParseDuration(fmt.Sprintf("%ds", respLongLived.Expires))
	dbUser.TokenExpires = time.Now().Add(duration)

	c.userDB.Save(dbUser)

	// Get User Profile info for determining tieStrength
	cacheProfile := c.getUserProfile(dbUser)
	c.logger.Printf("Got Profile %s: Likes(%d) Movies(%d) Events(%d) Music(%d) InsipirationalPeople(%d) Education(%d) Family(%d) Languages(%d)",
		dbUser.ID,
		len(cacheProfile.Likes),
		len(cacheProfile.Movies),
		len(cacheProfile.Events),
		len(cacheProfile.Music),
		len(cacheProfile.InspirationalPeople),
		len(cacheProfile.Education),
		len(cacheProfile.Family),
		len(cacheProfile.Languages),
	)
	c.userRedis.SaveProfile(dbUser.ID, cacheProfile)

	// This Queue can determine clique and tie-strength
	c.friendPublisher.Publish(dbUser)

	elapsed := time.Since(start)
	c.logger.Printf("Processed AuthQueue for %s in %s", dbUser.ID, elapsed)
	d.Ack(false)
}

func getLongLivedToken(u *domain.AuthUser) *facebookResponseLongLived {
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

func (c *Consumer) getUserProfile(user *domain.DBUser) *domain.CacheUserProfile {

	cacheProfile := &domain.CacheUserProfile{}

	res, _ := http.Get(fmt.Sprintf(
		"https://graph.facebook.com/v2.9/me?access_token=%s"+
			"&fields=gender,"+
			"age_range,"+
			"hometown,"+
			"location,"+
			"education,"+
			"favorite_teams,"+
			"inspirational_people,"+
			"languages,"+
			"sports,"+
			"political,"+
			"religion,"+
			"work,"+
			"family.limit(500){id},"+
			"music.limit(500){id},"+
			"movies.limit(500){id},"+
			"likes.limit(500){id},"+
			"events.limit(500){id}",
		user.LongLivedToken,
	))

	responseUserProfile := responseUserProfile{}
	err := json.NewDecoder(res.Body).Decode(&responseUserProfile)
	defer res.Body.Close()
	if err != nil {
		log.Printf("Error: %s", err)
	}
	cacheProfile = &domain.CacheUserProfile{
		Gender:    responseUserProfile.Gender,
		AgeRange:  strconv.Itoa(responseUserProfile.AgeRange.Min),
		Hometown:  responseUserProfile.Hometown.ID,
		Location:  responseUserProfile.Location.ID,
		Political: responseUserProfile.Political,
		Religion:  responseUserProfile.Religion,
	}

	for _, responseEducation := range responseUserProfile.Education {
		cacheProfile.Education = append(cacheProfile.Education, responseEducation.School.ID)
	}

	for _, responseMedia := range responseUserProfile.FavoriteTeams {
		cacheProfile.FavouriteTeams = append(cacheProfile.FavouriteTeams, responseMedia.ID)
	}

	for _, responseMedia := range responseUserProfile.InspirationalPeople {
		cacheProfile.InspirationalPeople = append(cacheProfile.InspirationalPeople, responseMedia.ID)
	}

	for _, responseMedia := range responseUserProfile.Languages {
		cacheProfile.Languages = append(cacheProfile.Languages, responseMedia.ID)
	}

	for _, responseMedia := range responseUserProfile.Sports {
		cacheProfile.Sports = append(cacheProfile.Sports, responseMedia.ID)
	}

	for _, responseMedia := range responseUserProfile.Work {
		cacheProfile.Work = append(cacheProfile.Work, responseMedia.ID)
	}

	for _, responseMedia := range responseUserProfile.Family.Data {
		cacheProfile.Family = append(cacheProfile.Family, responseMedia.ID)
	}

	for _, responseMedia := range responseUserProfile.Music.Data {
		cacheProfile.Music = append(cacheProfile.Music, responseMedia.ID)
	}

	for _, responseMedia := range responseUserProfile.Movies.Data {
		cacheProfile.Movies = append(cacheProfile.Movies, responseMedia.ID)
	}

	for _, responseMedia := range responseUserProfile.Likes.Data {
		cacheProfile.Likes = append(cacheProfile.Likes, responseMedia.ID)
	}

	for _, responseMedia := range responseUserProfile.Events.Data {
		cacheProfile.Events = append(cacheProfile.Events, responseMedia.ID)
	}

	return cacheProfile
}

type responseUserProfile struct {
	ID                  string              `json:"id"`
	Gender              string              `json:"gender"` // 'male', 'female'
	AgeRange            responseAgeRange    `json:"age_range"`
	Hometown            genericResponse     `json:"hometown"`
	Location            genericResponse     `json:"location"`
	Education           []responseEducation `json:"education"`
	FavoriteTeams       []genericResponse   `json:"favorite_teams"`
	InspirationalPeople []genericResponse   `json:"inspirational_people"`
	Languages           []genericResponse   `json:"languages"`
	Sports              []genericResponse   `json:"sports"`
	Work                []genericResponse   `json:"work"`
	Family              responseMedia       `json:"family"`
	Music               responseMedia       `json:"music"`
	Movies              responseMedia       `json:"movies"`
	Likes               responseMedia       `json:"likes"`
	Events              responseMedia       `json:"events"`
	Political           string              `json:"political"`
	Religion            string              `json:"religion"`
}

type responseAgeRange struct {
	Min int `json:"min"`
}

type responseEducation struct {
	School responseEducationSchool `json:"school"`
}

type responseEducationSchool struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"` // 'High School', 'College', 'Graduate School'
}

type genericResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type responseMedia struct {
	Data   []genericResponse `json:"data"`
	Paging responsePaging    `json:"paging"`
}

type responsePaging struct {
	Cursors responseCursors `json:"cursors"`
}

type responseCursors struct {
	Before string `json:"before"`
	After  string `json:"after"`
}
