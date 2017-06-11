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

// New - Returns a new FacebookPhoto.
func (m APIManager) New() *FacebookPhoto {
	return &FacebookPhoto{}
}

// Save - Saves the model across storages
func (m APIManager) Save(p *FacebookPhoto) error {
	jsonPhoto, _ := json.Marshal(p)
	m.redis.Do(
		"SET",
		fmt.Sprintf("photo:%s", p.FacebookPhotoID),
		jsonPhoto,
	)
	m.cacheLogger.Printf("Saved photo:%s", p.FacebookPhotoID)

	return nil
}

// FindByID - Returns a FacebookPhoto given an Id
func (m APIManager) FindByID(facebookID string) (*FacebookPhoto, error) {
	photo := &FacebookPhoto{}

	photoJSON, _ := redis.Bytes(m.redis.Do(
		"GET",
		fmt.Sprintf("photo:%s", facebookID),
	))

	if photoJSON != nil {
		json.Unmarshal(photoJSON, photo)
		m.cacheLogger.Printf("Got photo:%s", photo.FacebookPhotoID)

		return photo, nil
	}

	m.cacheLogger.Printf("Could not find photo:%s", facebookID)
	return nil, errors.New("Not found")

}
