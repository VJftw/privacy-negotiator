package photo

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
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
	userManager  user.Managable
	syncQueue    *SyncQueue
}

// NewController - Returns a new controller for photos.
func NewController(
	controllerLogger *log.Logger,
	renderer *render.Render,
	photoManager Managable,
	userManager user.Managable,
	syncQueue *SyncQueue,
) *Controller {
	return &Controller{
		logger:       controllerLogger,
		render:       renderer,
		photoManager: photoManager,
		userManager:  userManager,
		syncQueue:    syncQueue,
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

	router.Handle("/v1/photos/{id:[0-9]+}", negroni.New(
		middlewares.NewJWT(c.render),
		negroni.Wrap(http.HandlerFunc(c.postPhotoHandler)),
	)).Methods("PUT")

	log.Println("Set up Photo controller.")

}

func (c Controller) getPhotosHandler(w http.ResponseWriter, r *http.Request) {
	facebookUserID := middlewares.FBUserIDFromContext(r.Context())
	facebookUser, _ := c.userManager.FindByID(facebookUserID)
	idsJSON := r.URL.Query().Get("ids")
	var ids []string
	// JSON unmarshal url query ids
	json.Unmarshal([]byte(idsJSON), &ids)

	returnPhotos := []*FacebookPhoto{}
	// Find batch fb photo ids on redis.
	for _, facebookPhotoID := range ids {
		facebookPhoto, err := c.photoManager.FindByID(facebookPhotoID, facebookUser)
		if err == nil {
			returnPhotos = append(returnPhotos, facebookPhoto)
		}
	}

	c.render.JSON(w, http.StatusOK, returnPhotos)
}

func (c Controller) postPhotoHandler(w http.ResponseWriter, r *http.Request) {
	facebookUserID := middlewares.FBUserIDFromContext(r.Context())
	facebookUser, _ := c.userManager.FindByID(facebookUserID)
	photo, err := FromRequest(r)
	if err != nil {
		c.render.JSON(w, http.StatusBadRequest, nil)
		return
	}

	c.photoManager.Save(photo, facebookUser)

	c.syncQueue.Publish(photo)

	c.render.JSON(w, http.StatusCreated, photo)
}

func (c Controller) putPhotoHandler(w http.ResponseWriter, r *http.Request) {
	facebookUserID := middlewares.FBUserIDFromContext(r.Context())
	facebookUser, _ := c.userManager.FindByID(facebookUserID)
	vars := mux.Vars(r)
	photoID := vars["id"]

	photo, _ := c.photoManager.FindByID(photoID, facebookUser)

	FromPutRequest(r, photo, facebookUser)

	c.photoManager.Save(photo, facebookUser)

	c.render.JSON(w, http.StatusOK, photo)
}
