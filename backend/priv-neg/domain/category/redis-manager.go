package category

import (
	"log"

	"fmt"

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
func (m *RedisManager) Save(c string) error {
	redisConn := m.redis.Get()
	defer redisConn.Close()
	redisConn.Do(
		"SADD",
		"categories",
		c,
	)
	m.cacheLogger.Printf("Added category %s", c)

	return nil
}

// GetAll - Returns all of the categories
func (m *RedisManager) GetAll() []string {
	redisConn := m.redis.Get()
	defer redisConn.Close()
	categories, _ := redis.Strings(redisConn.Do(
		"SMEMBERS",
		"categories",
	))

	m.cacheLogger.Printf("Got categories")

	return categories
}

// GetCategoriesForUser - Returns user defined categories
func (m *RedisManager) GetCategoriesForUser(cacheUser *domain.CacheUser) []string {
	redisConn := m.redis.Get()
	defer redisConn.Close()
	categories, _ := redis.Strings(redisConn.Do(
		"SMEMBERS",
		fmt.Sprintf("u%s:categories", cacheUser.ID),
	))

	m.cacheLogger.Printf("Got categories for %s", cacheUser.ID)

	return categories
}

// AddCategoryForUser - Adds a user defined category
func (m *RedisManager) AddCategoryForUser(cacheUser *domain.CacheUser, c string) error {
	redisConn := m.redis.Get()
	defer redisConn.Close()
	redisConn.Do(
		"SADD",
		fmt.Sprintf("u%s:categories", cacheUser.ID),
		c,
	)
	m.cacheLogger.Printf("Added category %s to user %s", c, cacheUser.ID)

	return nil
}
