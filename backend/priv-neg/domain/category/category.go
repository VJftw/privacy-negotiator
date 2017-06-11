package category

// Category - The Category entity.
type Category struct {
	ID             string `json:"id" gorm:"primary_key"`
	Name           string `json:"name"`
	FacebookUserID string `json:"-" gorm:"ForeignKey:FacebookUserRefer"`
}
