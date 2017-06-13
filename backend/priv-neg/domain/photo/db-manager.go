package photo

import (
	"errors"
	"log"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
	"github.com/jinzhu/gorm"
)

// DBManager - Manages Photo entities on the Database.
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

// Save - Saves a given photo to the Database.
func (m *DBManager) Save(p *domain.DBPhoto) error {
	err := m.gorm.Where(domain.DBPhoto{ID: p.ID}).Assign(p).FirstOrCreate(p).Error
	if err != nil {
		return err
	}
	m.dbLogger.Printf("Saved photo %s", p.ID)

	return nil
}

// FindByID - Returns a Photo given its ID, nil if not found.
func (m *DBManager) FindByID(id string) (*domain.DBPhoto, error) {
	dbPhoto := &domain.DBPhoto{}

	m.gorm.Where("id = ?", id).First(dbPhoto)
	if dbPhoto.ID == "" {
		m.dbLogger.Printf("Could not find photo %s", id)
		return nil, errors.New("Not found")
	}

	m.dbLogger.Printf("Got photo %s", dbPhoto.ID)
	return dbPhoto, nil
}
