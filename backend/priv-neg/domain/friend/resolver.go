package friend

import (
	"net/http"

	"encoding/json"
	"errors"

	"fmt"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/utils"
)

// FromRequest - Returns a WebFriendship from a request and authenticated user.
func FromRequest(r *http.Request) (*domain.WebFriendship, error) {
	webFriendship := &domain.WebFriendship{}

	requestFriendship := &friendshipRequest{}
	err := json.NewDecoder(r.Body).Decode(requestFriendship)

	if err != nil {
		return nil, err
	}

	if requestFriendship.ID == "" {
		return nil, errors.New("missing id")
	}

	webFriendship.ID = requestFriendship.ID

	return webFriendship, nil
}

type friendshipRequest struct {
	ID string `json:"id"`
}

// FromPutRequest - Returns a modified CacheClique using a request
func FromPutRequest(r *http.Request, clique *domain.CacheClique, categories []string, userCategories []string) (*domain.CacheClique, error) {
	var jsonClique requestPUT
	err := json.NewDecoder(r.Body).Decode(&jsonClique)
	if err != nil {
		return nil, err
	}

	if len(jsonClique.Name) < 1 {
		return nil, errors.New("Name is too short")
	}
	clique.Name = jsonClique.Name

	clique.Categories = []string{}
	clique.UserCategories = []string{}
	for _, blockedCat := range jsonClique.BlockedCategories {
		found := false
		if utils.IsIn(blockedCat, categories) {
			found = true
			clique.Categories = append(clique.Categories, blockedCat)
		}
		if !found {
			if utils.IsIn(blockedCat, userCategories) {
				found = true
				clique.UserCategories = append(clique.UserCategories, blockedCat)
			}
		}
		if !found {
			return nil, fmt.Errorf("could not find category: %s", blockedCat)
		}
	}

	return clique, nil
}

type requestPUT struct {
	Name              string   `json:"name"`
	BlockedCategories []string `json:"categories"`
}
