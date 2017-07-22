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

// Save - Saves a given Friendship for a user ID to the Cache.
func (m *RedisManager) Save(uID string, f *domain.CacheFriendship) error {
	jsonFriendship, err := json.Marshal(f)
	if err != nil {
		return err
	}

	redisConn := m.redis.Get()
	defer redisConn.Close()
	redisConn.Do(
		"HSET",
		fmt.Sprintf("u%s:friends", uID),
		f.ID,
		jsonFriendship,
	)
	m.cacheLogger.Printf("Saved friendship %s:%s", uID, f.ID)

	return nil
}

// FindByIDAndUser - Returns a CacheFriendship given a friend user id and the authenticated user.
func (m *RedisManager) FindByIDAndUser(fUserID string, cacheUser *domain.CacheUser) (*domain.CacheFriendship, error) {
	cacheFriend := &domain.CacheFriendship{}

	redisConn := m.redis.Get()
	defer redisConn.Close()
	jsonFriendship, err := redis.Bytes(redisConn.Do(
		"HGET",
		fmt.Sprintf("u%s:friends", cacheUser.ID),
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

	m.cacheLogger.Printf("Could not find friendship %s:%s", cacheUser.ID, fUserID)
	return nil, errors.New("Not found")
}

// GetFriendIDsForAUserID - Returns friend IDs for a given user ID
func (m *RedisManager) GetFriendIDsForAUserID(fUserID string) []string {
	redisConn := m.redis.Get()
	defer redisConn.Close()
	friends, err := redis.Strings(redisConn.Do(
		"HKEYS",
		fmt.Sprintf("u%s:friends", fUserID),
	))

	if err != nil {
		m.cacheLogger.Printf("debug: %v", err)
		return friends
	}

	if friends != nil {
		m.cacheLogger.Printf("Got %d friends for %s", len(friends), fUserID)

		return friends
	}

	m.cacheLogger.Printf("Could not find friends for %s", fUserID)
	return friends
}

// GetCliqueIDsForAUserID - Returns clique IDs for a given user ID
func (m *RedisManager) GetCliqueIDsForAUserID(fUserID string) []string {
	redisConn := m.redis.Get()
	defer redisConn.Close()
	cliques, err := redis.Strings(redisConn.Do(
		"HKEYS",
		fmt.Sprintf("u%s:cliques", fUserID),
	))

	if err != nil {
		return cliques
	}

	if cliques != nil {
		m.cacheLogger.Printf("Got %d cliques for %s", len(cliques), fUserID)

		return cliques
	}

	m.cacheLogger.Printf("Could not find cliques for %s", fUserID)
	return cliques
}

// GetCliquesForAUserID - Returns cliques for a given user ID
func (m *RedisManager) GetCliquesForAUserID(fUserID string) []domain.CacheClique {
	redisConn := m.redis.Get()
	defer redisConn.Close()
	cliquesJSON, err := redis.ByteSlices(redisConn.Do(
		"HVALS",
		fmt.Sprintf("u%s:cliques", fUserID),
	))

	cliques := []domain.CacheClique{}
	if err != nil {
		m.cacheLogger.Printf("Error: %v", err)
		return cliques
	}

	if cliquesJSON != nil {
		for _, cliqueJSON := range cliquesJSON {
			clique := domain.CacheClique{}
			json.Unmarshal(cliqueJSON, &clique)
			cliques = append(cliques, clique)
		}
		m.cacheLogger.Printf("Got %d cliques for %s", len(cliques), fUserID)

		return cliques
	}

	m.cacheLogger.Printf("Could not find cliques for %s", fUserID)
	return cliques
}

// AddCliqueToUserID - Adds a Clique to a UserID
func (m *RedisManager) AddCliqueToUserID(fUserID string, clique *domain.CacheClique) {
	jsonClique, _ := json.Marshal(clique)
	redisConn := m.redis.Get()
	defer redisConn.Close()
	redisConn.Do(
		"HSET",
		fmt.Sprintf("u%s:cliques", fUserID),
		clique.ID,
		jsonClique,
	)

	m.cacheLogger.Printf("Added clique %s to %s", clique.ID, fUserID)
}

// RemoveCliqueByIDFromUserID - Removes a Clique from a UserID
func (m *RedisManager) RemoveCliqueByIDFromUserID(fUserID string, cliqueID string) {
	redisConn := m.redis.Get()
	defer redisConn.Close()
	redisConn.Do(
		"HDEL",
		fmt.Sprintf("u%s:cliques", fUserID),
		cliqueID,
	)

	m.cacheLogger.Printf("Removed clique %s from %s", cliqueID, fUserID)
}

// GetCliqueByIDAndUser - Returns a clique given its ID and user
func (m *RedisManager) GetCliqueByIDAndUser(cliqueID string, cacheUser *domain.CacheUser) (*domain.CacheClique, error) {
	cacheClique := &domain.CacheClique{}

	redisConn := m.redis.Get()
	defer redisConn.Close()
	jsonFriendship, err := redis.Bytes(redisConn.Do(
		"HGET",
		fmt.Sprintf("u%s:cliques", cacheUser.ID),
		cliqueID,
	))

	if err != nil {
		return nil, err
	}

	if jsonFriendship != nil {
		json.Unmarshal(jsonFriendship, cacheClique)
		m.cacheLogger.Printf("Got clique %s:%s", cacheUser.ID, cacheClique.ID)

		return cacheClique, nil
	}

	m.cacheLogger.Printf("Could not find clique %s:%s", cacheUser.ID, cliqueID)
	return nil, errors.New("Not found")
}
