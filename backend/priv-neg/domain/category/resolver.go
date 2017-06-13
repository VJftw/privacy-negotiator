package category

import (
	"encoding/json"
	"net/http"

	"errors"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
)

func FromRequest(r *http.Request, cacheUser *domain.CacheUser) (*domain.DBCategory, error) {
	dbCategory := &domain.DBCategory{}

	requestCategory := &categoryRequest{}
	err := json.NewDecoder(r.Body).Decode(requestCategory)
	if err != nil {
		return nil, err
	}

	if requestCategory.Name == "" {
		return nil, errors.New("Missing name.")
	}

	for _, existingCat := range cacheUser.Categories {
		if existingCat == requestCategory.Name {
			return nil, errors.New("Category already exists.")
		}
	}

	dbCategory.Name = requestCategory.Name
	dbCategory.UserID = cacheUser.ID

	return dbCategory, nil

}

type categoryRequest struct {
	Name string `json:"name"`
}
