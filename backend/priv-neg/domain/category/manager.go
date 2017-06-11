package category

// Managable - Manages Categories.
type Managable interface {
	New() *Category
	Save(*Category) error
	FindByID(string) (*Category, error)
}
