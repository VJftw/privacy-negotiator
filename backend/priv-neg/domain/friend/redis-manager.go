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

// FindByIDAndUser - Returns a CacheFriendship given a friend user id and the authenticated user.
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

// GetFriendIDsForAUserID - Returns friend IDs for a given user ID
func (m *RedisManager) GetFriendIDsForAUserID(fUserID string) []string {
	redisConn := m.redis.Get()
	defer redisConn.Close()
	jsonFriends, err := redis.Bytes(redisConn.Do(
		"HKEYS",
		fmt.Sprintf("%s:friends", fUserID),
		fUserID,
	))

	friends := []string{}

	if err != nil {
		return friends
	}

	if jsonFriends != nil {
		json.Unmarshal(jsonFriends, friends)
		m.cacheLogger.Printf("Got friends for %s", fUserID)

		return friends
	}

	m.cacheLogger.Printf("Could not find friends for %s", fUserID)
	return friends
}

// GetCliqueIDsForAUserID - Returns clique IDs for a given user ID
func (m *RedisManager) GetCliqueIDsForAUserID(fUserID string) []string {
	redisConn := m.redis.Get()
	defer redisConn.Close()
	jsonCliques, err := redis.Bytes(redisConn.Do(
		"HKEYS",
		fmt.Sprintf("%s:cliques", fUserID),
		fUserID,
	))

	cliques := []string{}

	if err != nil {
		return cliques
	}

	if jsonCliques != nil {
		json.Unmarshal(jsonCliques, cliques)
		m.cacheLogger.Printf("Got cliques for %s", fUserID)

		return cliques
	}

	m.cacheLogger.Printf("Could not find cliques for %s", fUserID)
	return cliques
}

// AddCliqueToUserID - Adds a given clique to a user ID
func (m *RedisManager) AddCliqueToUserID(fUserID string, clique *domain.CacheClique) {
	redisConn := m.redis.Get()
	defer redisConn.Close()
	redisConn.Do(
		"HSET",
		fmt.Sprintf("%s:cliques", fUserID),
		clique.ID,
		clique,
	)

	m.cacheLogger.Printf("Added clique %s to %s", clique.ID, fUserID)
}
