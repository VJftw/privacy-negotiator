package user

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type Manager interface {
	New() *FacebookUser
	Save(*FacebookUser) error
	GetInto(*FacebookUser, interface{}, ...interface{})
	FindByFacebookID(string) (*FacebookUser, error)
}

type UserManager struct {
	gorm *gorm.DB
}

func NewManager(gormDB *gorm.DB) Manager {
	return &UserManager{
		gorm: gormDB,
	}
}

func (m UserManager) New() *FacebookUser {
	return &FacebookUser{}
}

// Save - Saves the model across storages
func (m UserManager) Save(u *FacebookUser) error {
	m.gorm.Save(u)
	return nil
}

// GetInto - Searches the storages for a model identified by the query and places it into the given model reference.
// Returns true if found, false otherwise
func (m UserManager) GetInto(u *FacebookUser, query interface{}, args ...interface{}) {
	// check database
	m.gorm.Where(query, args...).First(u)
}

func (m UserManager) FindByFacebookID(uuid string) (*FacebookUser, error) {
	user := FacebookUser{}

	m.GetInto(&user, "uuid = ?", uuid)

	if len(user.FacebookUserID) < 1 {
		return nil, errors.New("Not found")
	}

	return &user, nil
}
