package user

// Managable - Manages FacebookUsers.
type Managable interface {
	New() *FacebookUser
	Save(*FacebookUser) error
	FindByFacebookID(string) (*FacebookUser, error)
}
