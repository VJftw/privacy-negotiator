package user

import (
	"errors"
	"log"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
	"github.com/jinzhu/gorm"
)

// DBManager - Manages User entities on the Database.
type DBManager struct {
	dbLogger *log.Logger
	gorm     *gorm.DB
}

// NewDBManager - Returns a new DBManager.
func NewDBManager(
	dbLogger *log.Logger,
	gorm *gorm.DB,
) *DBManager {
	return &DBManager{
		dbLogger: dbLogger,
		gorm:     gorm,
	}
}

// Save - Saves a given user to the Database.
func (m *DBManager) Save(u *domain.DBUser) error {
	err := m.gorm.Where(domain.DBUser{ID: u.ID}).Assign(u).FirstOrCreate(u).Error
	if err != nil {
		return err
	}
	m.dbLogger.Printf("Saved user %s", u.ID)

	return nil
}

// FindByID - Returns a user given its ID, nil if not found.
func (m *DBManager) FindByID(id string) (*domain.DBUser, error) {
	dbUserCliques := []domain.DBUserClique{}
	dbUser := &domain.DBUser{}

	err := m.gorm.Where("id = ?", id).First(dbUser).Error
	if err != nil {
		m.dbLogger.Printf("Error: %v", err)
		return nil, err
	}
	err = m.gorm.Model(dbUser).Related(&dbUserCliques, "DBUserCliques").Error
	if err != nil {
		m.dbLogger.Printf("Error: %v", err)
		return nil, err
	}
	dbUser.DBUserCliques = dbUserCliques

	if dbUser.ID == "" {
		m.dbLogger.Printf("Could not find user %s", id)
		return nil, errors.New("Not found")
	}

	m.dbLogger.Printf("Got user %v", dbUser)
	return dbUser, nil
}
