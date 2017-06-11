package category

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/middlewares"
)

// FromRequest - Returns a new Category using information from a request.
func FromRequest(r *http.Request) (*Category, error) {
	category := Category{}

	var rJSON map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&rJSON)
	if err != nil {
		return nil, err
	}

	if _, ok := rJSON["name"]; !ok {
		return nil, errors.New("Missing name")
	}

	fbUserID := middlewares.FBUserIDFromContext(r.Context())

	category.Name = rJSON["name"].(string)
	category.FacebookUserID = fbUserID

	return &category, nil
}
