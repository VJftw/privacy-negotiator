package auth

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	fb "github.com/huandu/facebook"
	"github.com/unrolled/render"
)

// Controller - Handles authentication
type Controller struct {
	render       *render.Render
	AuthResolver Resolver `inject:"auth.resolver"`
}

func (c Controller) Setup(router *mux.Router, renderer *render.Render) {
	c.render = renderer

	router.
		HandleFunc("/v1/auth", c.authHandler).
		Methods("POST")
}

func (c Controller) authHandler(w http.ResponseWriter, r *http.Request) {
	fbAuth := &FacebookAuth{}

	err := c.AuthResolver.FromRequest(fbAuth, r.Body)
	if err != nil {
		c.render.JSON(w, http.StatusBadRequest, nil)
		return
	}

	res, _ := fb.Get("/me", fb.Params{
		"fields":       "id",
		"access_token": fbAuth.AccessToken,
	})
	fmt.Println("here is my facebook first name:", res["id"])

	c.render.JSON(w, http.StatusUnauthorized, nil)
}
