package survey

import (
	"log"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
	"github.com/jinzhu/gorm"
)

// DBManager - Manages Survey entities on the Database.
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

// Save - Saves a given survey to the Database.
func (m *DBManager) Save(s *domain.DBSurvey) error {

	existingDBSurvey := domain.DBSurvey{}
	err := m.gorm.Debug().Where(domain.DBSurvey{
		UserID:  s.UserID,
		PhotoID: s.PhotoID,
	}).First(&existingDBSurvey).Error
	if err != nil { // Not found, create
		err = m.gorm.Debug().Create(s).Error
		if err != nil {
			return err
		}
	} else {
		m.gorm.Save(s)
	}

	m.dbLogger.Printf("Saved survey for %s", s.UserID)

	return nil
}
