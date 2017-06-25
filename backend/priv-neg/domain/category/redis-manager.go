package category

import (
	"fmt"
	"log"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
	"github.com/garyburd/redigo/redis"
)

// RedisManager - Manages user's categories on the Cache.
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

// Save - Saves a given Category to the Cache.
func (m *RedisManager) Save(u *domain.CacheUser, c string) error {
	redisConn := m.redis.Get()
	defer redisConn.Close()
	redisConn.Do(
		"SADD",
		fmt.Sprintf("%s:categories", u.ID),
		c,
	)
	m.cacheLogger.Printf("Saved %s:categories", u.ID)

	return nil
}

// FindByUser - Returns Categories given a user, nil if not found.
func (m *RedisManager) FindByUser(u *domain.CacheUser) ([]string, error) {
	redisConn := m.redis.Get()
	defer redisConn.Close()
	categories, err := redis.Strings(redisConn.Do(
		"SMEMBERS",
		fmt.Sprintf("%s:categories", u.ID),
	))
	if err != nil {
		m.cacheLogger.Printf("Could not find %s:categories", u.ID)
		return nil, err
	}

	m.cacheLogger.Printf("Found %s:categories", u.ID)
	return categories, nil
}
