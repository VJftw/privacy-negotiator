package friend

import (
	"log"

	"errors"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
	"github.com/jinzhu/gorm"
)

// DBManager - Manages Clique entities on the Database.
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

// Save - Saves a given Clique to the Database.
func (m *DBManager) Save(u *domain.DBClique) error {
	err := m.gorm.Where(domain.DBClique{ID: u.ID}).Assign(u).FirstOrCreate(u).Error
	if err != nil {
		return err
	}
	m.dbLogger.Printf("Saved clique %s", u.ID)

	return nil
}

// SaveUserClique - Saves a DBUserClique
func (m *DBManager) SaveUserClique(uC *domain.DBUserClique) error {

	existingDBUserClique := domain.DBUserClique{}

	err := m.gorm.Debug().Where(domain.DBUserClique{
		CliqueID: uC.CliqueID,
		UserID:   uC.UserID,
	}).First(&existingDBUserClique).Error
	if err != nil { // Not found, create
		err = m.gorm.Debug().Create(uC).Error
		if err != nil {
			return err
		}
	} else {
		categories := uC.Categories
		m.gorm.Debug().Model(uC).Association("Categories").Clear()
		uC.Categories = categories
		m.gorm.Debug().Save(uC)
	}

	m.dbLogger.Printf("Saved clique %s for user %s", uC.CliqueID, uC.UserID)

	return nil
}

// FindCliqueByID - Returns a clique given its ID, nil if not found.
func (m *DBManager) FindCliqueByID(id string) (*domain.DBClique, error) {
	dbUserCliques := []domain.DBUserClique{}
	dbClique := domain.DBClique{}

	err := m.gorm.Where("id = ?", id).First(&dbClique).Error
	if err != nil {
		m.dbLogger.Printf("Error: %v", err)
		return nil, err
	}

	err = m.gorm.Model(dbClique).Related(&dbUserCliques, "DBUserCliques").Error
	if err != nil {
		m.dbLogger.Printf("Error: %v", err)
		return nil, err
	}
	dbClique.DBUserCliques = dbUserCliques

	if dbClique.ID == "" {
		m.dbLogger.Printf("Could not find clique %s", id)
		return nil, errors.New("Not found")
	}

	m.dbLogger.Printf("Got clique %s", id)
	m.dbLogger.Printf("DEBUG: Clique Users: %v", dbClique.GetUserIDs())
	return &dbClique, nil
}

// GetUserCliquesByUser - Returns the UserCliques for a User.
func (m *DBManager) GetUserCliquesByUser(u domain.DBUser) ([]domain.DBUserClique, error) {

	dbUserCliques := []domain.DBUserClique{}

	err := m.gorm.Debug().Model(&u).Related(&dbUserCliques, "DBUserCliques").Error
	if err != nil {
		m.dbLogger.Printf("Error: %v", err)
		return nil, err
	}

	returnDbUserCliques := []domain.DBUserClique{}

	for _, dbUserClique := range dbUserCliques {
		categories := []domain.DBCategory{}
		err = m.gorm.Debug().Model(&dbUserClique).Related(&categories, "Categories").Error
		if err != nil {
			m.dbLogger.Printf("Error: %v", err)
			return nil, err
		}
		dbUserClique.Categories = categories
		returnDbUserCliques = append(returnDbUserCliques, dbUserClique)
	}

	return returnDbUserCliques, nil
}
