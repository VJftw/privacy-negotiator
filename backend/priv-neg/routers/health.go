package routers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

// Controller - Handles authentication
type Controller struct {
	logger *log.Logger
	render *render.Render
}

// NewHealthController - returns a new Controller for Health checks.
func NewHealthController(
	controllerLogger *log.Logger,
	renderer *render.Render,
) *Controller {
	return &Controller{
		logger: controllerLogger,
		render: renderer,
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
	c.render.JSON(w, http.StatusOK, nil)
}
