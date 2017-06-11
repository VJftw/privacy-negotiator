package photo

// FacebookPhoto - The FacebookPhoto entity.
type FacebookPhoto struct {
	FacebookPhotoID string   `json:"id" gorm:"primary_key"`
	Uploader        string   `json:"uploader"`
	TaggedUsers     []string `json:"taggedUsers"`

	Pending bool `json:"pending"`
}
