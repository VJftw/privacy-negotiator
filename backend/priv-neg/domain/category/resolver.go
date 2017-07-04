package category

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
)

// FromRequest - Returns a DBCategory from a httpRequest.
func FromRequest(r *http.Request, cacheUser *domain.CacheUser, userCategories []string) (*domain.DBCategory, error) {
	dbCategory := &domain.DBCategory{}

	requestCategory := &requestCategory{}
	err := json.NewDecoder(r.Body).Decode(requestCategory)
	if err != nil {
		return nil, err
	}

	if requestCategory.Name == "" {
		return nil, errors.New("missing name")
	}

	for _, existingCat := range userCategories {
		if existingCat == requestCategory.Name {
			return nil, errors.New("category already exists")
		}
	}

	dbCategory.Name = requestCategory.Name
	dbCategory.UserID = cacheUser.ID

	return dbCategory, nil

}

type requestCategory struct {
	Name string `json:"name"`
}
