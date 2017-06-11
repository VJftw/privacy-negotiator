package photo

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

// New - Returns a new FacebookPhoto.
func (m APIManager) New() *FacebookPhoto {
	return &FacebookPhoto{}
}

// Save - Saves the model across storages
func (m APIManager) Save(u *FacebookPhoto) error {
	jsonUser, _ := json.Marshal(u)
	m.Redis.Do(
		"SET",
		fmt.Sprintf("photo:%s", u.FacebookPhotoID),
		jsonUser,
	)
	m.CacheLogger.Printf("Saved photo:%s", u.FacebookPhotoID)

	// Should add to Queue

	return nil
}

// FindByID - Returns a FacebookPhoto given an Id
func (m APIManager) FindByID(facebookID string) (*FacebookPhoto, error) {
	user := &FacebookPhoto{}

	userJSON, _ := redis.Bytes(m.Redis.Do(
		"GET",
		fmt.Sprintf("photo:%s", facebookID),
	))

	if userJSON != nil {
		json.Unmarshal(userJSON, user)
		m.CacheLogger.Printf("Got photo:%s", user.FacebookPhotoID)

		return user, nil
	}

	m.CacheLogger.Printf("Could not find photo:%s", facebookID)
	return nil, errors.New("Not found")

}
