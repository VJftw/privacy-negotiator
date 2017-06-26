package category

import (
	"log"

	"errors"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
	"github.com/jinzhu/gorm"
)

// DBManager - Manages Category entities on the Database.
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
func (m *DBManager) Save(c *domain.DBCategory) error {
	err := m.gorm.Where(domain.DBCategory{
		Name: c.Name,
	}).Assign(c).FirstOrCreate(c).Error
	if err != nil {
		return err
	}
	m.dbLogger.Printf("Saved category %s", c.Name)

	return nil
}

// FindByName - Returns a DBCategory for a given name if it exists.
func (m *DBManager) FindByName(name string) (*domain.DBCategory, error) {
	dbCategory := &domain.DBCategory{}

	err := m.gorm.Where("name = ?", name).First(dbCategory).Error
	if err != nil {
		m.dbLogger.Printf("Error: %v", err)
		return nil, err
	}

	if dbCategory.Name == "" {
		m.dbLogger.Printf("Could not find category %s", name)
		return nil, errors.New("Not found")
	}

	return dbCategory, nil
}
