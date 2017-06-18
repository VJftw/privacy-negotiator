package friend

import (
	"log"

	"encoding/json"
	"fmt"

	"errors"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
	"github.com/garyburd/redigo/redis"
)

// RedisManager - Manages Friendship entities on the Cache.
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

// Save - Saves a given Friendship to the Cache.
func (m *RedisManager) Save(u *domain.CacheUser, f *domain.CacheFriendship) error {
	jsonFriendship, err := json.Marshal(f)
	if err != nil {
		return err
	}

	redisConn := m.redis.Get()
	defer redisConn.Close()
	redisConn.Do(
		"HSET",
		fmt.Sprintf("%s:friends", u.ID),
		f.ID,
		jsonFriendship,
	)
	m.cacheLogger.Printf("Saved friendship %s:%s", u.ID, f.ID)

	return nil
}

func (m *RedisManager) FindByIDAndUser(fUserID string, cacheUser *domain.CacheUser) (*domain.CacheFriendship, error) {
	cacheFriend := &domain.CacheFriendship{}

	redisConn := m.redis.Get()
	defer redisConn.Close()
	jsonFriendship, err := redis.Bytes(redisConn.Do(
		"HGET",
		fmt.Sprintf("%s:friends", cacheUser.ID),
		fUserID,
	))

	if err != nil {
		return nil, err
	}

	if jsonFriendship != nil {
		json.Unmarshal(jsonFriendship, cacheFriend)
		m.cacheLogger.Printf("Got friendship %s:%s", cacheUser.ID, cacheFriend.ID)

		return cacheFriend, nil
	}

	m.cacheLogger.Printf("Could not find friendship %s:%s", cacheUser.ID, cacheFriend.ID)
	return nil, errors.New("Not found")
}
