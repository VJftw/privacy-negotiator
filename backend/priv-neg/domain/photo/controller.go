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
	logger       *log.Logger
	render       *render.Render
	photoManager Managable
}

// NewController - Returns a new controller for photos.
func NewController(
	controllerLogger *log.Logger,
	renderer *render.Render,
	photoManager Managable,
) *Controller {
	return &Controller{
		logger:       controllerLogger,
		render:       renderer,
		photoManager: photoManager,
	}
}

// Setup - Sets up the routes for photos.
func (c Controller) Setup(router *mux.Router) {
	router.Handle("/v1/photos", negroni.New(
		middlewares.NewJWT(c.render),
		negroni.Wrap(http.HandlerFunc(c.getPhotosHandler)),
	)).Methods("GET")

	router.Handle("/v1/photos", negroni.New(
		middlewares.NewJWT(c.render),
		negroni.Wrap(http.HandlerFunc(c.postPhotoHandler)),
	)).Methods("POST")

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
		facebookPhoto, err := c.photoManager.FindByID(facebookPhotoID)
		if err == nil {
			returnPhotos = append(returnPhotos, facebookPhoto)
		}
	}

	c.render.JSON(w, http.StatusOK, returnPhotos)
}

func (c Controller) postPhotoHandler(w http.ResponseWriter, r *http.Request) {
	photo, err := FromRequest(r)
	if err != nil {
		c.render.JSON(w, http.StatusBadRequest, nil)
		return
	}

	c.photoManager.Save(photo)

	c.render.JSON(w, http.StatusCreated, photo)
}
