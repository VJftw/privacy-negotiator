package photo

// FacebookPhoto - The FacebookPhoto entity.
type FacebookPhoto struct {
	FacebookPhotoID string   `json:"id" gorm:"primary_key"`
	Pending         bool     `json:"pending"`
	TaggedUsers     []string `json:"taggedUsers"`
}
