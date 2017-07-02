package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
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
func (m *RedisManager) Save(u *domain.CacheUser) error {
	jsonUser, err := json.Marshal(u)
	if err != nil {
		return err
	}

	redisConn := m.redis.Get()
	defer redisConn.Close()
	redisConn.Do(
		"SET",
		fmt.Sprintf("u%s:info", u.ID),
		jsonUser,
	)
	m.cacheLogger.Printf("Saved user:%s", u.ID)

	return nil
}

// SaveProfile - Saves a profile to a user id
func (m *RedisManager) SaveProfile(id string, p *domain.CacheUserProfile) error {
	jsonUser, err := json.Marshal(p)
	if err != nil {
		return err
	}

	redisConn := m.redis.Get()
	defer redisConn.Close()
	redisConn.Do(
		"SET",
		fmt.Sprintf("u%s:profile", id),
		jsonUser,
	)
	m.cacheLogger.Printf("Saved user profile: %s", id)

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
func (m *RedisManager) FindByID(id string) (*domain.CacheUser, error) {
	user := &domain.CacheUser{}

	redisConn := m.redis.Get()
	defer redisConn.Close()
	userJSON, err := redis.Bytes(redisConn.Do(
		"GET",
		fmt.Sprintf("u%s:info", id),
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

// GetProfileByID - Returns a user Profile given an ID
func (m *RedisManager) GetProfileByID(id string) (*domain.CacheUserProfile, error) {
	cacheUserProfile := &domain.CacheUserProfile{}

	redisConn := m.redis.Get()
	defer redisConn.Close()
	userJSON, err := redis.Bytes(redisConn.Do(
		"GET",
		fmt.Sprintf("u%s:profile", id),
	))
	if err != nil {
		return nil, err
	}

	if userJSON != nil {
		json.Unmarshal(userJSON, cacheUserProfile)
		m.cacheLogger.Printf("Got u%s:profile", id)

		return cacheUserProfile, nil
	}

	m.cacheLogger.Printf("Could not find u%s:profile", id)
	return nil, errors.New("Not found")
}
