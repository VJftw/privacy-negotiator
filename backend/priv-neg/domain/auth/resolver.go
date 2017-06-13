package auth

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
)

// FromRequest - Validates and Transforms raw request data into a struct of Facebook credentials.
func FromRequest(b io.ReadCloser) (*user.WebUser, error) {
	user := &user.WebUser{}

	requestUser := &requestUser{}

	err := json.NewDecoder(b).Decode(requestUser)
	if err != nil {
		return nil, err
	}

	if requestUser.ID == "" {
		return nil, errors.New("Missing userID")
	}

	if requestUser.AccessToken == "" {
		return nil, errors.New("Missing accessToken")
	}

	user.ID = requestUser.ID
	user.ShortLivedToken = requestUser.AccessToken

	return user, nil
}

type requestUser struct {
	ID          string `json:"userID"`
	AccessToken string `json:"accessToken"`
}
