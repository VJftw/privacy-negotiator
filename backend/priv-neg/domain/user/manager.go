package user

// Managable - Manages FacebookUsers.
type Managable interface {
	Save(*User) error
	FindByID(string) (*FacebookUser, error)
}
