package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
)

// GraphAPI - Usage of the Facebook Graph API. This should be minimal in the API to reduce latency.
type GraphAPI interface {
	ValidateCredentials(*user.FacebookUser) bool
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

func (f *facebookGraphAPI) ValidateCredentials(fbAuth *user.FacebookUser) bool {
	res, err := http.Get(fmt.Sprintf("https://graph.facebook.com/v2.9/me?access_token=%s", fbAuth.ShortLivedToken))

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

	if me.ID == fbAuth.FacebookUserID {
		return true
	}

	return false
}
