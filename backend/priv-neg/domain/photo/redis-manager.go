package photo

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
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
func (m *RedisManager) Save(p *CachePhoto) error {
	jsonPhoto, err := json.Marshal(p)
	if err != nil {
		return err
	}

	redisConn := m.redis.Get()
	defer redisConn.Close()
	redisConn.Do(
		"SET",
		fmt.Sprintf("photo:%s", p.ID),
		jsonPhoto,
	)
	m.cacheLogger.Printf("Saved photo:%s", p.ID)

	return nil
}

// FindByID - Returns a Photo given its ID, nil if not found.
func (m *RedisManager) FindByID(id string) (*CachePhoto, error) {
	photo := &CachePhoto{}

	redisConn := m.redis.Get()
	defer redisConn.Close()
	jsonPhoto, err := redis.Bytes(redisConn.Do(
		"GET",
		fmt.Sprintf("photo:%s", id),
	))
	if err != nil {
		return nil, err
	}

	if jsonPhoto != nil {
		json.Unmarshal(jsonPhoto, photo)
		m.cacheLogger.Printf("Got photo:%s", photo.ID)

		return photo, nil
	}

	m.cacheLogger.Printf("Could not find photo:%s", photo.ID)
	return nil, errors.New("Not found")
}

// FindByIDWithUserCategories - Returns a WebPhoto with the user's chosen categories.
func (m *RedisManager) FindByIDWithUserCategories(id string, user *user.CacheUser) (*WebPhoto, error) {
	cachePhoto, err := m.FindByID(id)
	if err != nil {
		return nil, err
	}
	webPhoto := WebPhotoFromCachePhoto(cachePhoto)

	redisConn := m.redis.Get()
	defer redisConn.Close()
	jsonPhotoCategories, _ := redis.Bytes(redisConn.Do(
		"GET",
		fmt.Sprintf("%s:%s", webPhoto.ID, user.ID),
	))

	if jsonPhotoCategories != nil {
		json.Unmarshal(jsonPhotoCategories, webPhoto.Categories)
		m.cacheLogger.Printf("Got photo for user %s:%s", webPhoto.ID, user.ID)
	}

	return webPhoto, nil
}
