package photo

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/category"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/middlewares"
)

// FromRequest - Returns a Photo from a http request.
func FromRequest(r *http.Request) (*FacebookPhoto, error) {
	photo := &FacebookPhoto{}

	var photoJSON map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&photoJSON)
	if err != nil {
		return nil, err
	}

	if _, ok := photoJSON["id"]; !ok {
		return nil, errors.New("Missing id")
	}

	photo.FacebookPhotoID = photoJSON["id"].(string)
	photo.Uploader = middlewares.FBUserIDFromContext(r.Context())
	photo.Pending = true

	return photo, nil
}

// FromPutRequest - Modifies a given FacebookPhoto with data from a request for a given User.
func FromPutRequest(r *http.Request, p *FacebookPhoto, u *user.FacebookUser) error {
	var photoJSON requestPUT
	err := json.NewDecoder(r.Body).Decode(&photoJSON)
	if err != nil {
		return err
	}

	p.Categories = []category.Category{}

	for _, cat := range photoJSON.Categories {
		c := category.Category{
			Name:           cat,
			FacebookUserID: u.FacebookUserID,
		}

		p.Categories = append(p.Categories, c)
	}

	return nil

}

type requestPUT struct {
	Categories []string `json:"categories"`
}
