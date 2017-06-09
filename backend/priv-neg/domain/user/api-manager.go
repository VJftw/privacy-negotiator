package user

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

// New - Returns a new FacebookUser.
func (m APIManager) New() *FacebookUser {
	return &FacebookUser{}
}

// Save - Saves the model across storages
func (m APIManager) Save(u *FacebookUser) error {
	jsonUser, _ := json.Marshal(u)
	m.Redis.Do(
		"SET",
		fmt.Sprintf("user:%s", u.FacebookUserID),
		jsonUser,
	)
	m.CacheLogger.Printf("Saved user:%s\n", u.FacebookUserID)

	// Should add to Queue

	return nil
}

// FindByFacebookID - Returns a FacebookUser given an facebookUserId
func (m APIManager) FindByFacebookID(facebookID string) (*FacebookUser, error) {
	user := &FacebookUser{}

	userJSON, _ := m.Redis.Do(
		"GET",
		fmt.Sprintf("user:%s", facebookID),
	)

	if userJSON != nil {
		str, _ := userJSON.(string)
		json.Unmarshal([]byte(str), user)

		return user, nil
	}

	return nil, errors.New("Not found")

}
