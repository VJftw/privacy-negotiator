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

// New - Returns a new Category.
func (m WorkerManager) New() *Category {
	return &Category{}
}

// Save - Saves the model across storages
func (m WorkerManager) Save(u *Category) error {
	jsonUser, _ := json.Marshal(u)
	redisConn := m.redis.Get()
	defer redisConn.Close()
	redisConn.Do(
		"SET",
		fmt.Sprintf("category:%s", u.ID),
		jsonUser,
	)
	m.cacheLogger.Printf("Saved category:%s", u.ID)
	m.gorm.Where(
		Category{ID: u.ID},
	).Assign(u).FirstOrCreate(u)
	m.dbLogger.Printf("Saved category %s", u.ID)

	return nil
}

// FindByID - Returns a Category given an Id
func (m WorkerManager) FindByID(ID string) (*Category, error) {
	user := &Category{}

	// Check cache first.
	redisConn := m.redis.Get()
	defer redisConn.Close()
	userJSON, _ := redisConn.Do(
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