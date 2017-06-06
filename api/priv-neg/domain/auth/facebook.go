package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// GraphAPI - Usage of the Facebook Graph API. This should be minimal in the API to reduce latency.
type GraphAPI interface {
	ValidateCredentials(*FacebookAuth) bool
}

// NewGraphAPI - Returns the implementation of GraphAPI.
func NewGraphAPI() GraphAPI {
	return &facebookGraphAPI{}
}

type me struct {
	ID string `json:"id"`
}

type facebookGraphAPI struct {
}

func (f *facebookGraphAPI) ValidateCredentials(fbAuth *FacebookAuth) bool {
	res, err := http.Get(fmt.Sprintf("https://graph.facebook.com/v2.9/me?access_token=%s", fbAuth.AccessToken))

	if err != nil {
		log.Printf("Error: %s", err)
		return false
	}

	me := &me{}
	err = json.NewDecoder(res.Body).Decode(me)

	if err != nil {
		log.Printf("Error: %s", err)
		return false
	}

	if me.ID == fbAuth.UserID {
		return true
	}

	return false
}
