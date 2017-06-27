package friend

import (
	"net/http"

	"encoding/json"
	"errors"

	"fmt"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
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
func FromPutRequest(r *http.Request, clique *domain.CacheClique, categories []string) (*domain.CacheClique, error) {
	var jsonClique requestPUT
	err := json.NewDecoder(r.Body).Decode(&jsonClique)
	if err != nil {
		return nil, err
	}

	if len(jsonClique.Name) < 1 {
		return nil, errors.New("Name is too short")
	}
	clique.Name = jsonClique.Name

	for _, blockedCat := range jsonClique.BlockedCategories {
		found := false
		for _, cat := range categories {
			if cat == blockedCat {
				found = true
			}
		}
		if !found {
			return nil, fmt.Errorf("could not find category: %s", blockedCat)
		}
	}

	clique.Categories = jsonClique.BlockedCategories

	return clique, nil
}

type requestPUT struct {
	Name              string   `json:"name"`
	BlockedCategories []string `json:"blocked"`
}
