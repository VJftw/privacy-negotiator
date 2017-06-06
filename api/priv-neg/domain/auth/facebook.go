package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type GraphAPI interface {
	ValidateCredentials(*FacebookAuth) bool
}

type me struct {
	ID string `json:"id"`
}

func NewGraphAPI() GraphAPI {
	return &facebookGraphApi{}
}

type facebookGraphApi struct {
}

func (f *facebookGraphApi) ValidateCredentials(fbAuth *FacebookAuth) bool {
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
