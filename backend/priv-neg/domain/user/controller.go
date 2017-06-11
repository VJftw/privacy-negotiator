package user

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/middlewares"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

// Controller - Handles users
type Controller struct {
	render      *render.Render
	UserManager Managable `inject:"user.manager"`
}

// Setup - Sets up the Auth Controller
func (c Controller) Setup(router *mux.Router, renderer *render.Render) {
	c.render = renderer

	router.Handle("/v1/users", negroni.New(
		middlewares.NewJWT(renderer),
		negroni.Wrap(http.HandlerFunc(c.getUsersHandler)),
	)).Methods("GET")

	log.Println("Set up User controller.")

}

func (c Controller) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	idsJSON := r.URL.Query().Get("ids")
	var ids []string
	// JSON unmarshal url query ids
	json.Unmarshal([]byte(idsJSON), &ids)
	log.Println(ids)

	returnIds := []string{}
	// Find batch fb user ids on redis.
	for _, facebookUserID := range ids {
		_, err := c.UserManager.FindByFacebookID(facebookUserID)
		if err == nil {
			returnIds = append(returnIds, facebookUserID)
		}
	}

	c.render.JSON(w, http.StatusOK, returnIds)

}
