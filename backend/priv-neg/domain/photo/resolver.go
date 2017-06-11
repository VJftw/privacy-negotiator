package photo

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/middlewares"
)

func MultipleFromRequest(r *http.Request) (*[]FacebookPhoto, error) {
	photos := []FacebookPhoto{}

	var rJSON []map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&rJSON)
	if err != nil {
		return nil, err
	}

	fbUserID := middlewares.FBUserIDFromContext(r.Context())

	for _, photoJSON := range rJSON {
		if _, ok := photoJSON["id"]; ok {
			photo := FacebookPhoto{
				FacebookPhotoID: photoJSON["id"].(string),
				Uploader:        fbUserID,
				Pending:         true,
			}
			photos = append(photos, photo)
		}
	}

	return &photos, nil
}

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

	return photo, nil
}
