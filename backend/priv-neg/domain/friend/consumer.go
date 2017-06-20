package friend

import (
	"log"
	"time"

	"encoding/json"

	"fmt"
	"net/http"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/utils"
	"github.com/streadway/amqp"
)

// Consumer - Consumer for getting Community Detection.
type Consumer struct {
	queue     amqp.Queue
	channel   *amqp.Channel
	logger    *log.Logger
	userDB    *user.DBManager
	userRedis *user.RedisManager
}

// NewConsumer - Returns a new consumer.
func NewConsumer(
	queueLogger *log.Logger,
	ch *amqp.Channel,
	userDBManager *user.DBManager,
	userRedisManager *user.RedisManager,
) *Consumer {
	queue, err := ch.QueueDeclare(
		"community-detection", // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")
	return &Consumer{
		logger:    queueLogger,
		channel:   ch,
		queue:     queue,
		userDB:    userDBManager,
		userRedis: userRedisManager,
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

	webFriendship := &domain.WebFriendship{}
	json.Unmarshal(d.Body, webFriendship)

	fromCacheUser, _ := c.userRedis.FindByID(webFriendship.From)
	toCacheUser, _ := c.userRedis.FindByID(webFriendship.To)
	c.logger.Printf("Started processing for new friendship %s:%s", fromCacheUser.ID, toCacheUser.ID)

	// k=3
	// Given a new Facebook User.
	// * Get their friends.
	//   * Find which cliques their friends belong to.
	//   	* Can they join any of these cliques? (Are (k-1) friends in the same clique?)
	// 	 * Can we form any new cliques?
	// 		* (Will probably hit API limit) Use Mutual Friends from Graph API! if length of result is >= k then a new clique can be formed.
	//		* Use Redis Cache
	//			* <userID>:friends | [<userID>, <userID>]
	dbUser, _ := c.userDB.FindByID(toCacheUser.ID)
	friendIDs := c.getFacebookFriendsForUser(dbUser, "", nil)

	mutualCliqueIDs := c.getMutualCliqueIDsForUserIDs(friendIDs)

	// add user to mutualCliques

	// Form new Cliques
	for _, friendID := range friendIDs {
		mutualFriends := c.getMutualFriendIDsForUserIDs(dbUser.ID, friendID)
		if len(mutualFriends) >= 2 {
			// Form clique

		}
	}


	elapsed := time.Since(start)
	c.logger.Printf("Processed friendship %s:%s in %s", fromCacheUser.ID, toCacheUser.ID, elapsed)
}

func (c *Consumer) getMutualCliqueIDsForUserIDs(friendIDs []string) []string {
	mutualCliqueIDs := []string{}
	for i, friendID := range friendIDs {
		cliqueIDs := c.userRedis.GetCliqueIDsForUserID(friendID)
		for _, cliqueID := range cliqueIDs {
			for _, friendIDJ := range friendIDs[i:] {
				friendJCliqueIDs := c.userRedis.GetCliquesForUserID(friendIDJ)
				if isIn(cliqueID, friendJCliqueIDs) {
					mutualCliqueIDs = append(mutualCliqueIDs, cliqueID)
				}
			}
		}
	}

	return mutualCliqueIDs
}

func (c *Consumer) getMutualFriendIDsForUserIDs(a string, b string) []string {
	mutualFriends := []string{}
	aFriends := c.userRedis.getFriendsForUserID(a)
	bFriends := c.userRedis.getFriendsForUserID(b)

	for _, aFriend := range aFriends {
		for _, bFriend := range bFriends {
			if aFriend == bFriend {
				mutualFriends = append(mutualFriends, aFriend)
			}
		}
	}

	return mutualFriends
}

func isIn(needle string, haystack []string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}

func (c *Consumer) getFacebookFriendsForUser(user *domain.DBUser, offset string, friendIds []string) []string {

	if friendIds == nil {
		friendIds = []string{}
	}

	res, _ := http.Get(fmt.Sprintf(
		"https://graph.facebook.com/v2.9/%s/friends?access_token=%s&limit=500%s",
		user.ID,
		user.LongLivedToken,
		offset,
	))

	responseFriends := &responseFriends{}
	err := json.NewDecoder(res.Body).Decode(responseFriends)
	defer res.Body.Close()
	if err != nil {
		log.Printf("Error: %s", err)
	}

	for _, responseFriend := range responseFriends.Data {
		friendIds = append(friendIds, responseFriend.ID)
	}

	if responseFriends.Paging.Cursors.After != "" {
		return c.getFacebookFriendsForUser(
			user,
			fmt.Sprintf("&after=%s", responseFriends.Paging.Cursors.After),
			friendIds,
		)
	}

	return friendIds
}

type responseFriends struct {
	Data   []responseFriend `json:"data"`
	Paging responsePaging   `json:"paging"`
}

type responseFriend struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type responsePaging struct {
	Cursors responseCursors `json:"cursors"`
}

type responseCursors struct {
	Before string `json:"before"`
	After  string `json:"after"`
}
