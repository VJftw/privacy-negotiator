package auth

import (
	"encoding/json"
	"errors"
	"io"
)

type Resolver interface {
	FromRequest(*FacebookAuth, io.ReadCloser) error
}

type authResolver struct {
}

func NewResolver() Resolver {
	return &authResolver{}
}

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
