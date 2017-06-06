package auth

import (
	"encoding/json"
	"errors"
	"io"
)

// Resolver - Transformation into authentication structures.
type Resolver interface {
	FromRequest(*FacebookAuth, io.ReadCloser) error
}

type authResolver struct {
}

// NewResolver - Returns an implementation of Resolver.
func NewResolver() Resolver {
	return &authResolver{}
}

// FromRequest - Validates and Transforms raw request data into a struct of Facebook credentials.
func (a authResolver) FromRequest(fbAuth *FacebookAuth, b io.ReadCloser) error {
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

	fbAuth.AccessToken = rJSON["accessToken"].(string)
	fbAuth.UserID = rJSON["userID"].(string)

	return nil
}
