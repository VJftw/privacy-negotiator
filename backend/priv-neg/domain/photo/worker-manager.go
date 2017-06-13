package photo

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
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
func (m WorkerManager) Save(p *FacebookPhoto, u *user.FacebookUser) error {
	jsonPhoto, _ := json.Marshal(u)
	redisConn := m.redis.Get()
	defer redisConn.Close()
	redisConn.Do(
		"SET",
		fmt.Sprintf("photo:%s", p.FacebookPhotoID),
		jsonPhoto,
	)
	m.cacheLogger.Printf("Saved photo:%s", p.FacebookPhotoID)
	jsonWSMessage, _ := json.Marshal(websocket.Message{Type: "photo", Data: p})
	redisConn.Do(
		"PUBLISH",
		fmt.Sprintf("user:%s", p.Uploader),
		jsonWSMessage,
	)
	m.cacheLogger.Printf("Published photo %s to %s", p.FacebookPhotoID, p.Uploader)
	for user := range p.TaggedUsers {
		redisConn.Do(
			"PUBLISH",
			fmt.Sprintf("user:%v", user),
			jsonPhoto,
		)
		m.cacheLogger.Printf("Published photo %s to %v", p.FacebookPhotoID, user)
	}

	m.gorm.Where(
		FacebookPhoto{FacebookPhotoID: p.FacebookPhotoID},
	).Assign(p).FirstOrCreate(p)
	m.dbLogger.Printf("Saved photo %s", p.FacebookPhotoID)

	return nil
}

// FindByID - Returns a FacebookPhoto given an Id
func (m WorkerManager) FindByID(facebookID string, facebookUser *user.FacebookUser) (*FacebookPhoto, error) {
	photo := &FacebookPhoto{}

	redisConn := m.redis.Get()
	defer redisConn.Close()
	photoJSON, _ := redis.Bytes(redisConn.Do(
		"GET",
		fmt.Sprintf("photo:%s", facebookID),
	))

	if photoJSON != nil {
		json.Unmarshal(photoJSON, photo)
		m.cacheLogger.Printf("Got photo:%s", photo.FacebookPhotoID)

		photoCategoriesJSON, _ := redis.Bytes(redisConn.Do(
			"GET",
			fmt.Sprintf("%s:%s", photo.FacebookPhotoID, facebookUser.FacebookUserID),
		))

		if photoCategoriesJSON != nil {
			json.Unmarshal(photoCategoriesJSON, photo.Categories)
			m.cacheLogger.Printf("Got photo user %s:%s", photo.FacebookPhotoID, facebookUser.FacebookUserID)
		}

		return photo, nil
	}

	m.cacheLogger.Printf("Could not find photo:%s", facebookID)
	// return nil, errors.New("Not found")

	// Check DB. If in DB, update cache.
	// m.GetInto(user, "userId = ?", facebookID)
	//
	// if len(user.FacebookPhotoID) < 1 {
	// 	return nil, errors.New("Not found")
	// }

	return photo, nil
}
