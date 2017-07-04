package domain

// CachePhoto - The representation of a Photo stored on the Cache.
type CachePhoto struct {
	ID             string              `json:"id"`
	TaggedUsers    []string            `json:"taggedUsers"`
	Uploader       string              `json:"uploader"`
	Pending        bool                `json:"pending"`
	Categories     []string            `json:"categories"`
	UserCategories map[string][]string `json:"userCategories"`
	Conflicts      []CacheConflict     `json:"conflicts"`
	AllowedUserIDs []string            `json:"allowedUsers"`
	BlockedUserIDs []string            `json:"blockedUsers"`
}

// WebPhoto - The photo representation sent to a web client.
type WebPhoto struct {
	ID             string              `json:"id"`
	TaggedUsers    []string            `json:"taggedUsers"`
	Uploader       string              `json:"from"`
	Pending        bool                `json:"pending"`
	Categories     []string            `json:"categories"`
	UserCategories map[string][]string `json:"userCategories"`
	Conflicts      []CacheConflict     `json:"conflicts"`
	AllowedUserIDs []string            `json:"allowedUsers"`
	BlockedUserIDs []string            `json:"blockedUsers"`
}

// DBPhoto - The entity stored on the database
type DBPhoto struct {
	ID       string `gorm:"primary_key"`
	Uploader string

	TaggedUsers []DBUser `gorm:"many2many:photo_users"`

	Categories []DBCategory `gorm:"many2many:photo_categories"`
}

// TableName - Returns the table name for the entity.
func (p DBPhoto) TableName() string {
	return "photos"
}

// WebPhotoFromCachePhoto - Converts a CachePhoto into a WebPhoto.
func WebPhotoFromCachePhoto(cachePhoto *CachePhoto) *WebPhoto {
	if cachePhoto.TaggedUsers == nil {
		cachePhoto.TaggedUsers = []string{}
	}
	if cachePhoto.Categories == nil {
		cachePhoto.Categories = []string{}
	}
	if cachePhoto.AllowedUserIDs == nil {
		cachePhoto.AllowedUserIDs = []string{}
	}
	if cachePhoto.BlockedUserIDs == nil {
		cachePhoto.BlockedUserIDs = []string{}
	}
	if cachePhoto.Conflicts == nil {
		cachePhoto.Conflicts = []CacheConflict{}
	}
	if cachePhoto.UserCategories == nil {
		cachePhoto.UserCategories = map[string][]string{}
	}
	return &WebPhoto{
		ID:             cachePhoto.ID,
		TaggedUsers:    cachePhoto.TaggedUsers,
		Pending:        cachePhoto.Pending,
		Uploader:       cachePhoto.Uploader,
		Categories:     cachePhoto.Categories,
		UserCategories: cachePhoto.UserCategories,
		Conflicts:      cachePhoto.Conflicts,
		AllowedUserIDs: cachePhoto.AllowedUserIDs,
		BlockedUserIDs: cachePhoto.BlockedUserIDs,
	}
}

// CachePhotoFromWebPhoto - Converts a WebPhoto into a CachePhoto.
func CachePhotoFromWebPhoto(webPhoto *WebPhoto) *CachePhoto {
	if webPhoto.TaggedUsers == nil {
		webPhoto.TaggedUsers = []string{}
	}
	if webPhoto.Categories == nil {
		webPhoto.Categories = []string{}
	}
	if webPhoto.AllowedUserIDs == nil {
		webPhoto.AllowedUserIDs = []string{}
	}
	if webPhoto.BlockedUserIDs == nil {
		webPhoto.BlockedUserIDs = []string{}
	}
	if webPhoto.UserCategories == nil {
		webPhoto.UserCategories = map[string][]string{}
	}
	return &CachePhoto{
		ID:             webPhoto.ID,
		TaggedUsers:    webPhoto.TaggedUsers,
		Pending:        webPhoto.Pending,
		Uploader:       webPhoto.Uploader,
		Categories:     webPhoto.Categories,
		UserCategories: webPhoto.UserCategories,
	}
}

// DBPhotoFromCachePhoto - Converts a Cache Photo partially into a DBPhoto. Beware when Saving this.
func DBPhotoFromCachePhoto(cachePhoto *CachePhoto) *DBPhoto {
	dbPhoto := DBPhoto{
		ID:       cachePhoto.ID,
		Uploader: cachePhoto.Uploader,
	}

	for _, id := range cachePhoto.TaggedUsers {
		dbPhoto.TaggedUsers = append(dbPhoto.TaggedUsers, DBUser{ID: id})
	}

	for _, cat := range cachePhoto.Categories {
		dbPhoto.Categories = append(dbPhoto.Categories, DBCategory{Name: cat, UserID: "none"})
	}

	for userID, cats := range cachePhoto.UserCategories {
		for _, cat := range cats {
			dbPhoto.Categories = append(dbPhoto.Categories, DBCategory{Name: cat, UserID: userID})
		}
	}

	return &dbPhoto
}
