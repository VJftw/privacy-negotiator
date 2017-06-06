package user

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// Manager - Manages FacebookUsers.
type Manager interface {
	New() *FacebookUser
	Save(*FacebookUser) error
	GetInto(*FacebookUser, interface{}, ...interface{})
	FindByFacebookID(string) (*FacebookUser, error)
}

type userManager struct {
	gorm *gorm.DB
}

// NewManager - Returns an implementation of Manager.
func NewManager(gormDB *gorm.DB) Manager {
	return &userManager{
		gorm: gormDB,
	}
}

func (m userManager) New() *FacebookUser {
	return &FacebookUser{}
}

// Save - Saves the model across storages
func (m userManager) Save(u *FacebookUser) error {
	m.gorm.Save(u)
	return nil
}

// GetInto - Searches the storages for a model identified by the query and places it into the given model reference.
// Returns true if found, false otherwise
func (m userManager) GetInto(u *FacebookUser, query interface{}, args ...interface{}) {
	// check database
	m.gorm.Where(query, args...).First(u)
}

func (m userManager) FindByFacebookID(uuid string) (*FacebookUser, error) {
	user := FacebookUser{}

	m.GetInto(&user, "uuid = ?", uuid)

	if len(user.FacebookUserID) < 1 {
		return nil, errors.New("Not found")
	}

	return &user, nil
}
