package photo

import (
	"encoding/json"
	"errors"
	"net/http"

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

// FromPutRequest - Modifies a given FacebookPhoto with data from a request for a given User.
// func FromPutRequest(r *http.Request, p *FacebookPhoto, u *user.FacebookUser) error {
// 	var photoJSON requestPUT
// 	err := json.NewDecoder(r.Body).Decode(&photoJSON)
// 	if err != nil {
// 		return err
// 	}
//
// 	p.Categories = []category.Category{}
//
// 	for _, cat := range photoJSON.Categories {
// 		c := category.Category{
// 			Name:           cat,
// 			FacebookUserID: u.FacebookUserID,
// 		}
//
// 		p.Categories = append(p.Categories, c)
// 	}
//
// 	return nil
//
// }
//
// type requestPUT struct {
// 	Categories []string `json:"categories"`
// }
