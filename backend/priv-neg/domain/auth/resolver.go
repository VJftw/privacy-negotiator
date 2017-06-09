package auth

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
)

// Resolver - Transformation into authentication structures.
type Resolver interface {
	FromRequest(*user.FacebookUser, io.ReadCloser) error
}

type authResolver struct {
}

// NewResolver - Returns an implementation of Resolver.
func NewResolver() Resolver {
	return &authResolver{}
}

// FromRequest - Validates and Transforms raw request data into a struct of Facebook credentials.
func (a authResolver) FromRequest(fbAuth *user.FacebookUser, b io.ReadCloser) error {
	var rJSON map[string]interface{}

	err := json.NewDecoder(b).Decode(&rJSON)
	if err != nil {
		return err
	}

	if _, ok := rJSON["accessToken"]; !ok {
		return errors.New("Missing accessToken")
	}

	if _, ok := rJSON["userID"]; !ok {
		return errors.New("Missing userID")
	}

	fbAuth.ShortLivedToken = rJSON["accessToken"].(string)
	fbAuth.FacebookUserID = rJSON["userID"].(string)

	return nil
}
