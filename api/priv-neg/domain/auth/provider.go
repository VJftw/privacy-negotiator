package auth

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// Provider - For creating new Authentication Tokens.
type Provider interface {
	NewFromFacebookAuth(*FacebookAuth) *APIAuth
}

type authProvider struct {
}

// NewProvider - Returns an implementation of the Provider.
func NewProvider() Provider {
	return &authProvider{}
}

// NewFromFacebookAuth - Returns an AuthToken based off Facebook ID.
func (p authProvider) NewFromFacebookAuth(fbAuth *FacebookAuth) *APIAuth {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"fbUserID": fbAuth.UserID,
		"nbf":      time.Now().Unix(),
	})
	tokenString, _ := token.SignedString([]byte("hmacSecret"))

	return &APIAuth{
		Token: tokenString,
	}

}
