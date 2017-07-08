package survey

import (
	"log"

	"net/http"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/middlewares"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

// Controller - Handles Surveys.
type Controller struct {
	logger          *log.Logger
	render          *render.Render
	userRedis       *user.RedisManager
	surveyPublisher *Publisher
}

// NewController - Returns a new Controller for Surveys.
func NewController(
	controllerLogger *log.Logger,
	renderer *render.Render,
	userRedisManager *user.RedisManager,
	surveyPublisher *Publisher,
) *Controller {
	return &Controller{
		logger:          controllerLogger,
		render:          renderer,
		userRedis:       userRedisManager,
		surveyPublisher: surveyPublisher,
	}
}

// Setup - Sets up the routes for Surveys.
func (c Controller) Setup(router *mux.Router) {
	router.Handle("/v1/surveys", negroni.New(
		middlewares.NewJWT(c.render),
		negroni.Wrap(http.HandlerFunc(c.postSurveysHandler)),
	)).Methods("POST")
}

func (c Controller) postSurveysHandler(w http.ResponseWriter, r *http.Request) {
	facebookUserID := middlewares.FBUserIDFromContext(r.Context())
	cacheUser, _ := c.userRedis.FindByID(facebookUserID)

	dbSurvey, err := FromRequest(r, cacheUser)
	if err != nil {
		c.render.JSON(w, http.StatusBadRequest, nil)
		return
	}

	c.surveyPublisher.Publish(dbSurvey)

	c.render.JSON(w, http.StatusCreated, dbSurvey)
}
