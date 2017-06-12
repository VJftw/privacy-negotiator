package photo

import (
	"encoding/json"
	"errors"
	"net/http"

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
