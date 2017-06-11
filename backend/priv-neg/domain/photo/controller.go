package photo

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/middlewares"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

// Controller - Handles photos
type Controller struct {
	render       *render.Render
	PhotoManager Managable `inject:"photo.manager"`
}

func (c Controller) Setup(router *mux.Router, renderer *render.Render) {
	c.render = renderer

	router.Handle("/v1/photos", negroni.New(
		middlewares.NewJWT(renderer),
		negroni.Wrap(http.HandlerFunc(c.getPhotosHandler)),
	)).Methods("GET")

	log.Println("Set up Photo controller.")

}

func (c Controller) getPhotosHandler(w http.ResponseWriter, r *http.Request) {
	idsJSON := r.URL.Query().Get("ids")
	var ids []string
	// JSON unmarshal url query ids
	json.Unmarshal([]byte(idsJSON), &ids)
	log.Println(ids)

	returnPhotos := []*FacebookPhoto{}
	// Find batch fb photo ids on redis.
	for _, facebookPhotoID := range ids {
		facebookPhoto, err := c.PhotoManager.FindByID(facebookPhotoID)
		if err == nil {
			returnPhotos = append(returnPhotos, facebookPhoto)
		}
	}

	c.render.JSON(w, http.StatusOK, returnPhotos)

}

func (c Controller) putPhotosHandler(w http.ResponseWriter, r *http.Request) {
	fbUserID := middlewares.FBUserIDFromContext(r.Context())

}
