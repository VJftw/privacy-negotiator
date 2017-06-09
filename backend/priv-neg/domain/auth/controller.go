package auth

import (
	"net/http"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/queues"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

// Controller - Handles authentication
type Controller struct {
	render                         *render.Render
	AuthResolver                   Resolver                          `inject:"auth.resolver"`
	AuthProvider                   Provider                          `inject:"auth.provider"`
	GraphAPI                       GraphAPI                          `inject:"auth.graphAPI"`
	GetFacebookLongLivedTokenQueue *queues.GetFacebookLongLivedToken `inject:"queues.getFacebookLongLivedToken"`
}

// Setup - Sets up the Auth Controller
func (c Controller) Setup(router *mux.Router, renderer *render.Render) {
	c.render = renderer

	router.
		HandleFunc("/v1/auth", c.authHandler).
		Methods("POST")
}

func (c Controller) authHandler(w http.ResponseWriter, r *http.Request) {
	fbAuth := &user.FacebookUser{}

	err := c.AuthResolver.FromRequest(fbAuth, r.Body)
	if err != nil {
		c.render.JSON(w, http.StatusBadRequest, nil)
		return
	}

	if !c.GraphAPI.ValidateCredentials(fbAuth) {
		c.render.JSON(w, http.StatusUnauthorized, nil)
		return
	}

	token := c.AuthProvider.NewFromFacebookAuth(fbAuth)
	c.render.JSON(w, http.StatusCreated, token)

	// Add GetLongLivedToken to queue
	c.GetFacebookLongLivedTokenQueue.Publish(fbAuth)
}
