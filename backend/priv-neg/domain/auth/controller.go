package auth

import (
	"log"
	"net/http"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

// Controller - Handles authentication
type Controller struct {
	logger        *log.Logger
	render        *render.Render
	authPublisher *Publisher
	userRedis     *user.RedisManager
}

// NewController - returns a new Controller for Authentication.
func NewController(
	controllerLogger *log.Logger,
	renderer *render.Render,
	authPublisher *Publisher,
	userRedis *user.RedisManager,
) *Controller {
	return &Controller{
		logger:        controllerLogger,
		render:        renderer,
		authPublisher: authPublisher,
		userRedis:     userRedis,
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

	webUser, err := FromRequest(r.Body)
	if err != nil {
		c.render.JSON(w, http.StatusBadRequest, nil)
		return
	}

	if !ValidateFacebookCredentials(webUser) {
		c.render.JSON(w, http.StatusUnauthorized, nil)
		return
	}

	token := NewFromFacebookAuth(webUser)
	c.render.JSON(w, http.StatusCreated, token)

	// Add GetLongLivedToken to queue
	c.authPublisher.Publish(webUser)

	// Save to Redis
	cacheUser := domain.CacheUserFromAuthUser(webUser)
	c.userRedis.Save(cacheUser)
}
