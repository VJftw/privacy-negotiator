package category

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

// New - Returns a new Category.
func (m WorkerManager) New() *Category {
	return &Category{}
}

// Save - Saves the model across storages
func (m WorkerManager) Save(u *Category) error {
	jsonUser, _ := json.Marshal(u)
	m.Redis.Do(
		"SET",
		fmt.Sprintf("category:%s", u.ID),
		jsonUser,
	)
	m.CacheLogger.Printf("Saved category:%s", u.ID)
	m.Gorm.Save(u)
	m.DbLogger.Printf("Saved category %s", u.ID)

	return nil
}

// FindByID - Returns a Category given an Id
func (m WorkerManager) FindByID(ID string) (*Category, error) {
	user := &Category{}

	// Check cache first.
	userJSON, _ := m.Redis.Do(
		"GET",
		fmt.Sprintf("category:%s", ID),
	)

	if userJSON != nil {
		str, _ := userJSON.(string)
		json.Unmarshal([]byte(str), user)

		return user, nil
	}

	// Check DB. If in DB, update cache.
	// m.GetInto(user, "userId = ?", facebookID)
	//
	// if len(user.CategoryID) < 1 {
	// 	return nil, errors.New("Not found")
	// }

	return user, nil
}
