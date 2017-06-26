package photo

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
	"github.com/garyburd/redigo/redis"
)

// RedisManager - Manages Photo entities on the Cache.
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

// Save - Saves a given Photo to the Cache.
func (m *RedisManager) Save(p *domain.CachePhoto) error {
	jsonPhoto, err := json.Marshal(p)
	if err != nil {
		return err
	}

	redisConn := m.redis.Get()
	defer redisConn.Close()
	redisConn.Do(
		"SET",
		fmt.Sprintf("p%s:info", p.ID),
		jsonPhoto,
	)
	m.cacheLogger.Printf("Saved photo:%s", p.ID)

	return nil
}

// FindByID - Returns a Photo given its ID, nil if not found.
func (m *RedisManager) FindByID(id string) (*domain.CachePhoto, error) {
	photo := &domain.CachePhoto{}

	redisConn := m.redis.Get()
	defer redisConn.Close()
	jsonPhoto, err := redis.Bytes(redisConn.Do(
		"GET",
		fmt.Sprintf("p%s:info", id),
	))
	if err != nil {
		return nil, err
	}

	if jsonPhoto != nil {
		json.Unmarshal(jsonPhoto, photo)
		m.cacheLogger.Printf("Got photo:%s", photo.ID)

		photo.Categories = m.GetCategoriesForPhoto(photo)

		return photo, nil
	}

	m.cacheLogger.Printf("Could not find photo:%s", photo.ID)
	return nil, errors.New("Not found")
}

// GetCategoriesForPhoto - Returns all categories for a given photo
func (m *RedisManager) GetCategoriesForPhoto(p *domain.CachePhoto) []string {
	redisConn := m.redis.Get()
	defer redisConn.Close()
	categories, _ := redis.Strings(redisConn.Do(
		"SMEMBERS",
		fmt.Sprintf("p%s:categories", p.ID),
	))

	m.cacheLogger.Printf("Got categories for %s", p.ID)

	return categories
}

// SaveCategoryForPhoto - Adds a given category to the Photo
func (m *RedisManager) SaveCategoryForPhoto(p *domain.CachePhoto, c string) error {
	redisConn := m.redis.Get()
	defer redisConn.Close()
	redisConn.Do(
		"SADD",
		fmt.Sprintf("p%s:categories", p.ID),
		c,
	)

	m.cacheLogger.Printf("Added categories for %s", p.ID)

	return nil
}
