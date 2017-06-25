package routers

import (
	"log"
	"net/http"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/persisters"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

// Controller - Handles authentication
type Controller struct {
	logger     *log.Logger
	render     *render.Render
	publishers []persisters.Publisher
}

// NewHealthController - returns a new Controller for Health checks.
func NewHealthController(
	controllerLogger *log.Logger,
	renderer *render.Render,
	publishers []persisters.Publisher,
) *Controller {
	return &Controller{
		logger:     controllerLogger,
		render:     renderer,
		publishers: publishers,
	}
}

// Setup - Sets up the Auth Controller
func (c Controller) Setup(router *mux.Router) {
	router.
		HandleFunc("/v1/health", c.healthHandler).
		Methods("GET")
	log.Println("Set up Health controller.")
}

func (c Controller) healthHandler(w http.ResponseWriter, r *http.Request) {
	totalMessages := 0
	for _, publisher := range c.publishers {
		totalMessages = totalMessages + publisher.GetMessageTotal()
	}

	c.render.JSON(w, http.StatusOK, &health{
		MessageTotal: totalMessages,
	})
}

type health struct {
	MessageTotal int `json:"queueSize"`
}
