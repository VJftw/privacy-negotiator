package domain

// DBCategory - The Category entity.
type DBCategory struct {
	Name string `gorm:"primary_key"`

	User   DBUser
	UserID string `gorm:"primary_key"`

	Photos []DBPhoto `gorm:"many2many:photo_categories"`
}

// WebCategory - Category representation to web_app.
type WebCategory struct {
	Name string `json:"name"`
}