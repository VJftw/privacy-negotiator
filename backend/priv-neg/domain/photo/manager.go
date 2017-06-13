package photo

import "github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"

// Managable - Manages FacebookPhotos.
type Managable interface {
	New() *FacebookPhoto
	Save(*FacebookPhoto, *user.FacebookUser) error
	FindByID(string, *user.FacebookUser) (*FacebookPhoto, error)
}
