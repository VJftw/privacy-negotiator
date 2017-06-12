package photo

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/routers/websocket"
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

// New - Returns a new FacebookPhoto.
func (m WorkerManager) New() *FacebookPhoto {
	return &FacebookPhoto{}
}

// Save - Saves the model across storages
func (m WorkerManager) Save(u *FacebookPhoto) error {
	jsonUser, _ := json.Marshal(u)
	redisConn := m.redis.Get()
	defer redisConn.Close()
	redisConn.Do(
		"SET",
		fmt.Sprintf("photo:%s", u.FacebookPhotoID),
		jsonUser,
	)
	m.cacheLogger.Printf("Saved photo:%s", u.FacebookPhotoID)
	jsonWSMessage, _ := json.Marshal(websocket.Message{Type: "photo", Data: u})
	redisConn.Do(
		"PUBLISH",
		fmt.Sprintf("user:%s", u.Uploader),
		jsonWSMessage,
	)
	m.cacheLogger.Printf("Published photo %s to %s", u.FacebookPhotoID, u.Uploader)
	for user := range u.TaggedUsers {
		redisConn.Do(
			"PUBLISH",
			fmt.Sprintf("user:%v", user),
			jsonUser,
		)
		m.cacheLogger.Printf("Published photo %s to %v", u.FacebookPhotoID, user)
	}

	m.gorm.Where(
		FacebookPhoto{FacebookPhotoID: u.FacebookPhotoID},
	).Assign(u).FirstOrCreate(u)
	m.dbLogger.Printf("Saved photo %s", u.FacebookPhotoID)

	return nil
}

// FindByID - Returns a FacebookPhoto given an Id
func (m WorkerManager) FindByID(facebookID string) (*FacebookPhoto, error) {
	user := &FacebookPhoto{}

	// Check cache first.
	redisConn := m.redis.Get()
	defer redisConn.Close()
	userJSON, _ := redisConn.Do(
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
