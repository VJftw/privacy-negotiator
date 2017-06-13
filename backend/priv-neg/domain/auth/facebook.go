package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
)

type me struct {
	ID string `json:"id"`
}

// ValidateFacebookCredentials - Validates given Facebook credentials with the Graph API.
func ValidateFacebookCredentials(webUser *domain.AuthUser) bool {
	res, err := http.Get(fmt.Sprintf("https://graph.facebook.com/v2.9/me?access_token=%s", webUser.ShortLivedToken))

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

	if me.ID == webUser.ID {
		return true
	}

	return false
}
