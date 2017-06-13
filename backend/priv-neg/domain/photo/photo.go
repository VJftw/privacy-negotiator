package photo

import (
	"os/user"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/category"
	"github.com/lib/pq"
)

// FacebookPhoto - The FacebookPhoto entity.
type FacebookPhoto struct {
	FacebookPhotoID string         `json:"id" gorm:"primary_key"`
	Uploader        string         `json:"uploader"`
	TaggedUsers     pq.StringArray `json:"taggedUsers" gorm:"type:varchar(64)[]"`

	Pending bool `json:"pending" gorm:"-"`

	Categories []category.Category `json:"categories" gorm:"many2many:photo_categories"`
}

// CachePhoto - The representation of a Photo stored on the Cache.
type CachePhoto struct {
	ID          string   `json:"id"`
	TaggedUsers []string `json:"taggedUsers"`
	Uploader    string   `json:"uploader"`
}

// WebPhoto - The photo representation sent to a web client.
type WebPhoto struct {
	ID          string        `json:"id"`
	TaggedUsers []string      `json:"taggedUsers"`
	Uploader    *user.WebUser `json:"from"`
	Categories  []string      `json:"categories"`
}

// Photo - The entity stored on the database
type Photo struct {
	ID       string `gorm:"primary_key"`
	Uploader string

	TaggedUsers []user.User `gorm:"many2many:photo_users"`

	Categories []category.Category `gorm:"many2many:photo_categories"`
}
