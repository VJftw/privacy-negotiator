package photo

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

// New - Returns a new FacebookPhoto.
func (m WorkerManager) New() *FacebookPhoto {
	return &FacebookPhoto{}
}

// Save - Saves the model across storages
func (m WorkerManager) Save(u *FacebookPhoto) error {
	jsonUser, _ := json.Marshal(u)
	m.Redis.Do(
		"SET",
		fmt.Sprintf("photo:%s", u.FacebookPhotoID),
		jsonUser,
	)
	m.CacheLogger.Printf("Saved photo:%s", u.FacebookPhotoID)
	m.Gorm.Save(u)
	m.DbLogger.Printf("Saved photo %s", u.FacebookPhotoID)

	return nil
}

// FindByID - Returns a FacebookPhoto given an Id
func (m WorkerManager) FindByID(facebookID string) (*FacebookPhoto, error) {
	user := &FacebookPhoto{}

	// Check cache first.
	userJSON, _ := m.Redis.Do(
		"GET",
		fmt.Sprintf("photo:%s", facebookID),
	)

	if userJSON != nil {
		str, _ := userJSON.(string)
		json.Unmarshal([]byte(str), user)

		return user, nil
	}

	// Check DB. If in DB, update cache.
	// m.GetInto(user, "userId = ?", facebookID)
	//
	// if len(user.FacebookPhotoID) < 1 {
	// 	return nil, errors.New("Not found")
	// }

	return user, nil
}
