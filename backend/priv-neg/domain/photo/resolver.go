package photo

import (
	"encoding/json"
	"errors"
	"net/http"

	"fmt"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
)

// FromRequest - Returns a Photo from a http request.
func FromRequest(r *http.Request, user *domain.CacheUser) (*domain.WebPhoto, error) {
	photo := &domain.WebPhoto{}

	requestPhoto := &photoRequest{}
	err := json.NewDecoder(r.Body).Decode(requestPhoto)
	if err != nil {
		return nil, err
	}

	if requestPhoto.ID == "" {
		return nil, errors.New("Missing id")
	}

	photo.ID = requestPhoto.ID
	photo.Uploader = user.ID
	photo.Pending = true

	return photo, nil
}

type photoRequest struct {
	ID         string   `json:"id"`
	Categories []string `json:"categories"`
}

//FromPutRequest - Modifies a given FacebookPhoto with data from a request for a given User.
func FromPutRequest(r *http.Request, p *domain.CachePhoto, u *domain.CacheUser) (*domain.WebPhoto, error) {
	var jsonPhoto requestPUT
	err := json.NewDecoder(r.Body).Decode(&jsonPhoto)
	if err != nil {
		return nil, err
	}

	for _, cat := range jsonPhoto.Categories {
		found := false
		for _, existCat := range u.Categories {
			if cat == existCat {
				found = true
			}
		}
		if !found {
			return nil, fmt.Errorf("could not find category: %s", cat)
		}
	}

	webPhoto := &domain.WebPhoto{
		ID: p.ID,
	}

	webPhoto.Categories = []string{}

	for _, cat := range jsonPhoto.Categories {
		webPhoto.Categories = append(webPhoto.Categories, cat)
	}

	return webPhoto, nil

}

type requestPUT struct {
	Categories []string `json:"categories"`
}
