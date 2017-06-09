package auth

import (
	"time"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
	jwt "github.com/dgrijalva/jwt-go"
)

// Provider - For creating new Authentication Tokens.
type Provider interface {
	NewFromFacebookAuth(*user.FacebookUser) *APIAuth
}

type authProvider struct {
}

// NewProvider - Returns an implementation of the Provider.
func NewProvider() Provider {
	return &authProvider{}
}

// NewFromFacebookAuth - Returns an AuthToken based off Facebook ID.
func (p authProvider) NewFromFacebookAuth(fbAuth *user.FacebookUser) *APIAuth {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"fbUserID": fbAuth.FacebookUserID,
		"nbf":      time.Now().Unix(),
	})
	tokenString, _ := token.SignedString([]byte("hmacSecret"))

	return &APIAuth{
		Token: tokenString,
	}

}
