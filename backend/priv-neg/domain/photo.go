package domain

// CachePhoto - The representation of a Photo stored on the Cache.
type CachePhoto struct {
	ID          string   `json:"id"`
	TaggedUsers []string `json:"taggedUsers"`
	Uploader    string   `json:"uploader"`
	Pending     bool     `json:"pending"`
}

// WebPhoto - The photo representation sent to a web client.
type WebPhoto struct {
	ID          string   `json:"id"`
	TaggedUsers []string `json:"taggedUsers"`
	Uploader    string   `json:"from"`
	Pending     bool     `json:"pending"`
	Categories  []string `json:"categories"`
}

// DBPhoto - The entity stored on the database
type DBPhoto struct {
	ID       string `gorm:"primary_key"`
	Uploader string

	TaggedUsers []DBUser `gorm:"many2many:photo_users"`

	Categories []DBCategory `gorm:"many2many:photo_categories"`
}

// WebPhotoFromCachePhoto - Converts a CachePhoto into a WebPhoto.
func WebPhotoFromCachePhoto(cachePhoto *CachePhoto) *WebPhoto {
	return &WebPhoto{
		ID:          cachePhoto.ID,
		TaggedUsers: cachePhoto.TaggedUsers,
		Pending:     cachePhoto.Pending,
		Uploader:    cachePhoto.Uploader,
		Categories:  []string{},
	}
}

// CachePhotoFromWebPhoto - Converts a WebPhoto into a CachePhoto.
func CachePhotoFromWebPhoto(webPhoto *WebPhoto) *CachePhoto {
	return &CachePhoto{
		ID:          webPhoto.ID,
		TaggedUsers: webPhoto.TaggedUsers,
		Pending:     webPhoto.Pending,
		Uploader:    webPhoto.Uploader,
	}
}
