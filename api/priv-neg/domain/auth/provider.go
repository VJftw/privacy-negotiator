package auth

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type Provider interface {
	NewFromFacebookAuth(*FacebookAuth) *ApiAuth
}

type authProvider struct {
}

func NewProvider() Provider {
	return &authProvider{}
}

func (p authProvider) NewFromFacebookAuth(fbAuth *FacebookAuth) *ApiAuth {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"fbUserID": fbAuth.UserID,
		"nbf":      time.Now().Unix(),
	})
	tokenString, _ := token.SignedString([]byte("hmacSecret"))

	return &ApiAuth{
		Token: tokenString,
	}

}
