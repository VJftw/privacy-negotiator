package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
)

// APIManager - Implementation of Managable.
type APIManager struct {
	redis       *redis.Pool
	cacheLogger *log.Logger
}

// NewAPIManager - Returns an implementation of Manager.
func NewAPIManager(cacheLogger *log.Logger, redis *redis.Pool) Managable {
	return &APIManager{
		redis:       redis,
		cacheLogger: cacheLogger,
	}
}

// New - Returns a new FacebookUser.
func (m APIManager) New() *FacebookUser {
	return &FacebookUser{}
}

// Save - Saves the model across storages
func (m APIManager) Save(u *FacebookUser) error {
	jsonUser, _ := json.Marshal(u)
	redisConn := m.redis.Get()
	defer redisConn.Close()
	redisConn.Do(
		"SET",
		fmt.Sprintf("user:%s", u.FacebookUserID),
		jsonUser,
	)
	m.cacheLogger.Printf("Saved user:%s", u.FacebookUserID)

	// Should add to Queue

	return nil
}

// FindByID - Returns a FacebookUser given an facebookUserId
func (m APIManager) FindByID(facebookID string) (*FacebookUser, error) {
	user := &FacebookUser{}

	redisConn := m.redis.Get()
	defer redisConn.Close()
	userJSON, _ := redis.Bytes(redisConn.Do(
		"GET",
		fmt.Sprintf("user:%s", facebookID),
	))

	if userJSON != nil {
		json.Unmarshal(userJSON, user)
		m.cacheLogger.Printf("Got user:%s", user.FacebookUserID)

		return user, nil
	}

	m.cacheLogger.Printf("Could not find user:%s", facebookID)
	return nil, errors.New("Not found")

}
