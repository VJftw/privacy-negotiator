package auth

import (
	"log"
	"net/http"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

// Controller - Handles authentication
type Controller struct {
	logger      *log.Logger
	render      *render.Render
	authQueue   *LongAuthQueue
	userManager user.Managable
}

// NewController - returns a new Controller for Authentication.
func NewController(
	controllerLogger *log.Logger,
	renderer *render.Render,
	authQueue *LongAuthQueue,
	userManager user.Managable,
) *Controller {
	return &Controller{
		logger:      controllerLogger,
		render:      renderer,
		authQueue:   authQueue,
		userManager: userManager,
	}
}

// Setup - Sets up the Auth Controller
func (c Controller) Setup(router *mux.Router) {
	router.
		HandleFunc("/v1/auth", c.authHandler).
		Methods("POST")
	log.Println("Set up Auth controller.")
}

func (c Controller) authHandler(w http.ResponseWriter, r *http.Request) {
	fbUser := &user.FacebookUser{}

	err := FromRequest(fbUser, r.Body)
	if err != nil {
		c.render.JSON(w, http.StatusBadRequest, nil)
		return
	}

	if !ValidateFacebookCredentials(fbUser) {
		c.render.JSON(w, http.StatusUnauthorized, nil)
		return
	}

	token := NewFromFacebookAuth(fbUser)
	c.render.JSON(w, http.StatusCreated, token)

	// Add GetLongLivedToken to queue
	c.authQueue.Publish(fbUser)

	// Save to Redis
	c.userManager.Save(fbUser)
}
