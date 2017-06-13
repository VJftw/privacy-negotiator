package category

import (
	"log"

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
		Name:   c.Name,
		UserID: c.UserID,
	}).Assign(c).FirstOrCreate(c).Error
	if err != nil {
		return err
	}
	m.dbLogger.Printf("Saved category %s:%s", c.UserID, c.Name)

	return nil
}
