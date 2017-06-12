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
	dbLogger    *log.Logger
	gorm        *gorm.DB
	cacheLogger *log.Logger
	redis       *redis.Pool
}

// NewWorkerManager - Returns an implementation of Manager.
func NewWorkerManager(
	dbLogger *log.Logger,
	gorm *gorm.DB,
	cacheLogger *log.Logger,
	redis *redis.Pool,
) Managable {
	return &WorkerManager{
		dbLogger:    dbLogger,
		gorm:        gorm,
		cacheLogger: cacheLogger,
		redis:       redis,
	}
}

// New - Returns a new FacebookUser.
func (m WorkerManager) New() *FacebookUser {
	return &FacebookUser{}
}

// Save - Saves the model across storages
func (m WorkerManager) Save(u *FacebookUser) error {
	jsonUser, _ := json.Marshal(u)
	redisConn := m.redis.Get()
	defer redisConn.Close()
	redisConn.Do(
		"SET",
		fmt.Sprintf("user:%s", u.FacebookUserID),
		jsonUser,
	)
	m.cacheLogger.Printf("Saved user:%s", u.FacebookUserID)
	m.gorm.Save(u)
	m.dbLogger.Printf("Saved user %s", u.FacebookUserID)

	return nil
}

// FindByID - Returns a FacebookUser given an facebookUserId
func (m WorkerManager) FindByID(facebookID string) (*FacebookUser, error) {
	user := &FacebookUser{}

	// Check cache first.
	redisConn := m.redis.Get()
	defer redisConn.Close()
	userJSON, _ := redisConn.Do(
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
