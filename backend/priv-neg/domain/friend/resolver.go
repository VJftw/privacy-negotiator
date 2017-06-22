package friend

import (
	"net/http"

	"encoding/json"
	"errors"

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
