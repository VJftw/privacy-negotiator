package friend

import (
	"net/http"

	"encoding/json"
	"errors"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
)

func FromRequest(r *http.Request, user *domain.CacheUser) (*domain.WebFriendship, error) {
	webFriendship := &domain.WebFriendship{
		From: user.ID,
	}

	requestFriendship := &friendshipRequest{}
	err := json.NewDecoder(r.Body).Decode(requestFriendship)

	if err != nil {
		return nil, err
	}

	if requestFriendship.ID == "" {
		return nil, errors.New("missing id")
	}

	webFriendship.To = requestFriendship.ID

	return webFriendship, nil
}

type friendshipRequest struct {
	ID string `json:"id"`
}
