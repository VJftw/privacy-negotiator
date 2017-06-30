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

	existingDBPhoto := domain.DBPhoto{}
	err := m.gorm.Debug().Where("id = ?", p.ID).First(&existingDBPhoto).Error
	if err != nil { // Not found, create
		err = m.gorm.Debug().Create(p).Error
		if err != nil {
			return err
		}
	} else {
		categories := p.Categories
		m.gorm.Debug().Model(p).Association("Categories").Clear()
		p.Categories = categories
		m.gorm.Debug().Save(p)
	}

	m.dbLogger.Printf("Saved photo %s", p.ID)

	return nil
}

// FindByID - Returns a Photo given its ID, nil if not found.
func (m *DBManager) FindByID(id string) (*domain.DBPhoto, error) {
	dbPhoto := &domain.DBPhoto{}
	dbCategories := []domain.DBCategory{}
	dbTaggedUsers := []domain.DBUser{}

	m.gorm.Where(
		"id = ?", id,
	).First(
		dbPhoto,
	).Related(
		&dbCategories,
		"Categories",
	).Related(
		&dbTaggedUsers,
		"TaggedUsers",
	)
	dbPhoto.Categories = dbCategories
	dbPhoto.TaggedUsers = dbTaggedUsers

	if dbPhoto.ID == "" {
		m.dbLogger.Printf("Could not find photo %s", id)
		return nil, errors.New("Not found")
	}

	m.dbLogger.Printf("Got photo %s", dbPhoto.ID)
	return dbPhoto, nil
}
