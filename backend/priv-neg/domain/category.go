package domain

// DBCategory - The Category entity.
type DBCategory struct {
	Name        string         `gorm:"primary_key"`
	Photos      []DBPhoto      `gorm:"many2many:photo_categories"`
	UserCliques []DBUserClique `gorm:"many2many:user_clique_categories"`
}

func (c DBCategory) TableName() string {
	return "categories"
}

// WebCategory - Category representation to web_app.
type WebCategory struct {
	Name string `json:"name"`
}
