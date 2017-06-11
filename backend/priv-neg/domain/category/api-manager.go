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
	Redis       redis.Conn  `inject:"persister.cache"`
	CacheLogger *log.Logger `inject:"logger.cache"`
}

// NewAPIManager - Returns an implementation of Manager.
func NewAPIManager() Managable {
	return &APIManager{}
}

// New - Returns a new Category.
func (m APIManager) New() *Category {
	return &Category{}
}

// Save - Saves the model across storages
func (m APIManager) Save(c *Category) error {
	jsonUser, _ := json.Marshal(c)
	m.Redis.Do(
		"SET",
		fmt.Sprintf("category:%s", c.ID),
		jsonUser,
	)
	m.CacheLogger.Printf("Saved category:%s", c.ID)

	// Should add to Queue

	return nil
}

// FindByID - Returns a Category given an Id
func (m APIManager) FindByID(ID string) (*Category, error) {
	user := &Category{}

	userJSON, _ := redis.Bytes(m.Redis.Do(
		"GET",
		fmt.Sprintf("category:%s", ID),
	))

	if userJSON != nil {
		json.Unmarshal(userJSON, user)
		m.CacheLogger.Printf("Got category:%s", user.ID)

		return user, nil
	}

	m.CacheLogger.Printf("Could not find category:%s", ID)
	return nil, errors.New("Not found")

}
