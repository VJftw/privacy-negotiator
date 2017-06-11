package category

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
)

// APIManager - Implementation of Managable.
type APIManager struct {
	redis       redis.Conn
	cacheLogger *log.Logger
}

// NewAPIManager - Returns an implementation of Manager.
func NewAPIManager(cacheLogger *log.Logger, redis redis.Conn) Managable {
	return &APIManager{
		redis:       redis,
		cacheLogger: cacheLogger,
	}
}

// New - Returns a new Category.
func (m APIManager) New() *Category {
	return &Category{}
}

// Save - Saves the model across storages
func (m APIManager) Save(c *Category) error {
	jsonUser, _ := json.Marshal(c)
	m.redis.Do(
		"SET",
		fmt.Sprintf("category:%s", c.ID),
		jsonUser,
	)
	m.cacheLogger.Printf("Saved category:%s", c.ID)

	// Should add to Queue

	return nil
}

// FindByID - Returns a Category given an Id
func (m APIManager) FindByID(ID string) (*Category, error) {
	user := &Category{}

	userJSON, _ := redis.Bytes(m.redis.Do(
		"GET",
		fmt.Sprintf("category:%s", ID),
	))

	if userJSON != nil {
		json.Unmarshal(userJSON, user)
		m.cacheLogger.Printf("Got category:%s", user.ID)

		return user, nil
	}

	m.cacheLogger.Printf("Could not find category:%s", ID)
	return nil, errors.New("Not found")

}
