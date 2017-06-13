package category

import (
	"log"
	"net/http"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/middlewares"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

// Handles Categories.
type Controller struct {
	logger            *log.Logger
	render            *render.Render
	userRedis         *user.RedisManager
	categoryPublisher *Publisher
}

// NewController - Returns a new Controller for Categories.
func NewController(
	controllerLogger *log.Logger,
	renderer *render.Render,
	userRedisManager *user.RedisManager,
	categoryPublisher *Publisher,
) *Controller {
	return &Controller{
		logger:            controllerLogger,
		render:            renderer,
		userRedis:         userRedisManager,
		categoryPublisher: categoryPublisher,
	}
}

// Setup - Sets up the routes for Categories.
func (c Controller) Setup(router *mux.Router) {
	router.Handle("/v1/categories", negroni.New(
		middlewares.NewJWT(c.render),
		negroni.Wrap(http.HandlerFunc(c.getCategoriesHandler)),
	)).Methods("GET")
}

func (c Controller) getCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	facebookUserID := middlewares.FBUserIDFromContext(r.Context())
	cacheUser, _ := c.userRedis.FindByID(facebookUserID)

	c.render.JSON(w, http.StatusOK, cacheUser.Categories)
}

func (c Controller) postCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	facebookUserID := middlewares.FBUserIDFromContext(r.Context())
	cacheUser, _ := c.userRedis.FindByID(facebookUserID)

	dbCategory, err := FromRequest(r, cacheUser)
	if err != nil {
		c.render.JSON(w, http.StatusBadRequest, nil)
		return
	}

	cacheUser.Categories = append(cacheUser.Categories, dbCategory.Name)
	c.userRedis.Save(cacheUser)

	c.categoryPublisher.Publish(dbCategory)

	c.render.JSON(w, http.StatusOK, cacheUser.Categories)
}
