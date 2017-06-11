package user

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
)

// WorkerManager - Implementation of Managable.
type WorkerManager struct {
	Gorm        *gorm.DB    `inject:"persister.db"`
	DbLogger    *log.Logger `inject:"logger.db"`
	Redis       redis.Conn  `inject:"persister.cache"`
	CacheLogger *log.Logger `inject:"logger.cache"`
}

// NewWorkerManager - Returns an implementation of Manager.
func NewWorkerManager() Managable {
	return &WorkerManager{}
}

// New - Returns a new FacebookUser.
func (m WorkerManager) New() *FacebookUser {
	return &FacebookUser{}
}

// Save - Saves the model across storages
func (m WorkerManager) Save(u *FacebookUser) error {
	jsonUser, _ := json.Marshal(u)
	m.Redis.Do(
		"SET",
		fmt.Sprintf("user:%s", u.FacebookUserID),
		jsonUser,
	)
	m.CacheLogger.Printf("Saved user:%s", u.FacebookUserID)
	m.Gorm.Save(u)
	m.DbLogger.Printf("Saved user %s", u.FacebookUserID)

	return nil
}

// FindByFacebookID - Returns a FacebookUser given an facebookUserId
func (m WorkerManager) FindByFacebookID(facebookID string) (*FacebookUser, error) {
	user := &FacebookUser{}

	// Check cache first.
	userJSON, _ := m.Redis.Do(
		"GET",
		fmt.Sprintf("user:%s", facebookID),
	)

	if userJSON != nil {
		str, _ := userJSON.(string)
		json.Unmarshal([]byte(str), user)

		return user, nil
	}

	// Check DB. If in DB, update cache.
	// m.GetInto(user, "userId = ?", facebookID)
	//
	// if len(user.FacebookUserID) < 1 {
	// 	return nil, errors.New("Not found")
	// }

	return user, nil
}
