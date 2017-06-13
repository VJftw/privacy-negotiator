package auth

import (
	"time"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
	jwt "github.com/dgrijalva/jwt-go"
)

// NewFromFacebookAuth - Returns an AuthToken based off Facebook ID.
func NewFromFacebookAuth(webUser *user.AuthUser) *APIAuth {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"fbUserID": webUser.ID,
		"nbf":      time.Now().Unix(),
	})
	tokenString, _ := token.SignedString([]byte("hmacSecret"))

	return &APIAuth{
		Token: tokenString,
	}

}
