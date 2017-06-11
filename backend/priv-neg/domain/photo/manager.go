package photo

// Managable - Manages FacebookPhotos.
type Managable interface {
	New() *FacebookPhoto
	Save(*FacebookPhoto) error
	FindByID(string) (*FacebookPhoto, error)
}
