package photo

// Managable - Manages FacebookPhotos.
type Managable interface {
	New() *FacebookPhoto
	Save(*FacebookPhoto) error
	FindByFacebookID(string) (*FacebookPhoto, error)
}
