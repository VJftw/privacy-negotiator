package user

// Managable - Manages FacebookUsers.
type Managable interface {
	New() *FacebookUser
	Save(*FacebookUser) error
	FindByID(string) (*FacebookUser, error)
}
