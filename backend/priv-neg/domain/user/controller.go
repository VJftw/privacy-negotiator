package user

import (
	"fmt"
	"net/http"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/middlewares"
	"github.com/gorilla/context"
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
}

func (c Controller) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	facebookUserID := context.Get(r, "fbUserID")
	str, _ := facebookUserID.(string)
	facebookUser, _ := c.UserManager.FindByFacebookID(str)

	fmt.Println(facebookUser)

}
