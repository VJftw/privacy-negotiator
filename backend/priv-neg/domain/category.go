package domain

// DBCategory - The Category entity.
type DBCategory struct {
	Name        string         `gorm:"primary_key"`
	Photos      []DBPhoto      `gorm:"many2many:photo_categories"`
	UserCliques []DBUserClique `gorm:"many2many:user_clique_categories"`
	UserID      string         `gorm:"primary_key"` // For personal categories/context
}

// QueueCategory - Represents a category stored in the queue for persistence.
type QueueCategory struct {
	Name   string `json:"name"`
	UserID string `json:"userId"`
}

// TableName - Returns the table name for the entity.
func (c DBCategory) TableName() string {
	return "categories"
}

// WebCategory - Category representation to web_app.
type WebCategory struct {
	Name     string `json:"name"`
	Personal bool   `json:"personal"`
}

// QueueCategoryFromDBCategory - Returns a QueueCategory from a DBCategory
func QueueCategoryFromDBCategory(dbCategory *DBCategory) *QueueCategory {
	return &QueueCategory{
		Name:   dbCategory.Name,
		UserID: dbCategory.UserID,
	}
}

// DBCategoryFromQueueCategory - Returns a DBCategory from a QueueCategory
func DBCategoryFromQueueCategory(queueCategory *QueueCategory) *DBCategory {
	return &DBCategory{
		Name:   queueCategory.Name,
		UserID: queueCategory.UserID,
	}
}
