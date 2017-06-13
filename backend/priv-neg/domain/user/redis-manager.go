package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/routers/websocket"
	"github.com/garyburd/redigo/redis"
)

// RedisManager - Manages user entities on the Cache.
type RedisManager struct {
	redis       *redis.Pool
	cacheLogger *log.Logger
}

// NewRedisManager - Returns a new RedisManager.
func NewRedisManager(cacheLogger *log.Logger, redis *redis.Pool) *RedisManager {
	return &RedisManager{
		redis:       redis,
		cacheLogger: cacheLogger,
	}
}

// Save - Saves a given user to the Cache.
func (m *RedisManager) Save(u *CacheUser) error {
	jsonUser, err := json.Marshal(u)
	if err != nil {
		return err
	}

	redisConn := m.redis.Get()
	defer redisConn.Close()
	redisConn.Do(
		"SET",
		fmt.Sprintf("user:%s", u.ID),
		jsonUser,
	)
	m.cacheLogger.Printf("Saved user:%s", u.ID)

	return nil
}

// Publish - Publishes a given message to a user pubsub channel (websocket).
func (m *RedisManager) Publish(uID string, channel string, data interface{}) {
	redisConn := m.redis.Get()
	defer redisConn.Close()
	jsonWSMessage, _ := json.Marshal(websocket.Message{Type: channel, Data: data})
	redisConn.Do(
		"PUBLISH",
		fmt.Sprintf("user:%s", uID),
		jsonWSMessage,
	)
	m.cacheLogger.Printf("Published %s to user:%s", channel, uID)
}

// FindByID - Returns a user given its ID, nil if not found.
func (m *RedisManager) FindByID(id string) (*CacheUser, error) {
	user := &CacheUser{}

	redisConn := m.redis.Get()
	defer redisConn.Close()
	userJSON, err := redis.Bytes(redisConn.Do(
		"GET",
		fmt.Sprintf("user:%s", id),
	))
	if err != nil {
		return nil, err
	}

	if userJSON != nil {
		json.Unmarshal(userJSON, user)
		m.cacheLogger.Printf("Got user:%s", user.ID)

		return user, nil
	}

	m.cacheLogger.Printf("Could not find user:%s", user.ID)
	return nil, errors.New("Not found")
}
