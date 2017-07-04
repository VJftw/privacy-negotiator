package category

import (
	"log"
	"net/http"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/middlewares"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

// Controller - Handles Categories.
type Controller struct {
	logger            *log.Logger
	render            *render.Render
	userRedis         *user.RedisManager
	categoryRedis     *RedisManager
	categoryPublisher *Publisher
}

// NewController - Returns a new Controller for Categories.
func NewController(
	controllerLogger *log.Logger,
	renderer *render.Render,
	userRedisManager *user.RedisManager,
	categoryRedisManager *RedisManager,
	categoryPublisher *Publisher,
) *Controller {
	return &Controller{
		logger:            controllerLogger,
		render:            renderer,
		userRedis:         userRedisManager,
		categoryRedis:     categoryRedisManager,
		categoryPublisher: categoryPublisher,
	}
}

// Setup - Sets up the routes for Categories.
func (c Controller) Setup(router *mux.Router) {
	router.Handle("/v1/categories", negroni.New(
		middlewares.NewJWT(c.render),
		negroni.Wrap(http.HandlerFunc(c.getCategoriesHandler)),
	)).Methods("GET")

	router.Handle("/v1/categories", negroni.New(
		middlewares.NewJWT(c.render),
		negroni.Wrap(http.HandlerFunc(c.postCategoriesHandler)),
	)).Methods("POST")
}

func (c Controller) getCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	webCategories := []domain.WebCategory{}
	categories := c.categoryRedis.GetAll()

	for _, cat := range categories {
		webCategories = append(webCategories, domain.WebCategory{
			Name:     cat,
			Personal: false,
		})
	}

	// get User categories and append.
	userID := middlewares.FBUserIDFromContext(r.Context())
	cacheUser, _ := c.userRedis.FindByID(userID)
	userCategories := c.categoryRedis.GetCategoriesForUser(cacheUser)
	for _, cat := range userCategories {
		webCategories = append(webCategories, domain.WebCategory{
			Name:     cat,
			Personal: true,
		})
	}

	c.render.JSON(w, http.StatusOK, webCategories)
}

func (c Controller) postCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	facebookUserID := middlewares.FBUserIDFromContext(r.Context())
	cacheUser, _ := c.userRedis.FindByID(facebookUserID)

	globalCategories := c.categoryRedis.GetAll()
	categories := c.categoryRedis.GetCategoriesForUser(cacheUser)

	dbCategory, err := FromRequest(r, cacheUser, append(globalCategories, categories...))
	if err != nil {
		c.render.JSON(w, http.StatusBadRequest, nil)
		return
	}

	c.categoryRedis.AddCategoryForUser(cacheUser, dbCategory.Name)
	queueCategory := domain.QueueCategoryFromDBCategory(dbCategory)
	c.categoryPublisher.Publish(queueCategory)

	c.render.JSON(w, http.StatusCreated, append(categories, dbCategory.Name))
}
