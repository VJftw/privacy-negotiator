package photo

import "github.com/lib/pq"

// FacebookPhoto - The FacebookPhoto entity.
type FacebookPhoto struct {
	FacebookPhotoID string         `json:"id" gorm:"primary_key"`
	Uploader        string         `json:"uploader"`
	TaggedUsers     pq.StringArray `json:"taggedUsers" gorm:"type:varchar(64)[]"`

	Pending bool `json:"pending" gorm:"-"`
}
