package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
)

type me struct {
	ID string `json:"id"`
}

// ValidateFacebookCredentials - Validates given Facebook credentials with the Graph API.
func ValidateFacebookCredentials(fbAuth *user.FacebookUser) bool {
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
