package category

import (
	"log"
	"net/http"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/middlewares"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

// Controller - Handles photos
type Controller struct {
	render          *render.Render
	CategoryManager Managable `inject:"category.manager"`
}

// Setup - Sets up the Auth Controller
func (c Controller) Setup(router *mux.Router, renderer *render.Render) {
	c.render = renderer

	router.Handle("/v1/categories", negroni.New(
		middlewares.NewJWT(renderer),
		negroni.Wrap(http.HandlerFunc(c.getCategoriesHandler)),
	)).Methods("GET")

	router.Handle("/v1/categories", negroni.New(
		middlewares.NewJWT(renderer),
		negroni.Wrap(http.HandlerFunc(c.postCategoriesHandler)),
	)).Methods("POST")

	log.Println("Set up Category controller.")

}

func (c Controller) getCategoriesHandler(w http.ResponseWriter, r *http.Request) {

	c.render.JSON(w, http.StatusOK, nil)

}

func (c Controller) postCategoriesHandler(w http.ResponseWriter, r *http.Request) {

	category, err := FromRequest(r)
	if err != nil {
		c.render.JSON(w, http.StatusBadRequest, nil)
		return
	}

	c.CategoryManager.Save(category)

	c.render.JSON(w, http.StatusCreated, nil)
}
