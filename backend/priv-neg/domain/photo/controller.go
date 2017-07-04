package photo

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/category"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/middlewares"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/utils"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

// Controller - Handles photos
type Controller struct {
	logger                *log.Logger
	render                *render.Render
	photoRedis            *RedisManager
	userRedis             *user.RedisManager
	categoryRedis         *category.RedisManager
	photoPublisher        *Publisher
	photoPersistPublisher *PersistPublisher
}

// NewController - Returns a new controller for photos.
func NewController(
	controllerLogger *log.Logger,
	renderer *render.Render,
	photoRedisManager *RedisManager,
	userRedisManager *user.RedisManager,
	categoryRedisManager *category.RedisManager,
	photoPublisher *Publisher,
	photoPersistPublisher *PersistPublisher,
) *Controller {
	return &Controller{
		logger:                controllerLogger,
		render:                renderer,
		photoRedis:            photoRedisManager,
		userRedis:             userRedisManager,
		photoPublisher:        photoPublisher,
		categoryRedis:         categoryRedisManager,
		photoPersistPublisher: photoPersistPublisher,
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

	router.Handle("/v1/photos/{id}", negroni.New(
		middlewares.NewJWT(c.render),
		negroni.Wrap(http.HandlerFunc(c.putPhotoHandler)),
	)).Methods("PUT")

	log.Println("Set up Photo controller.")

}

func (c Controller) getPhotosHandler(w http.ResponseWriter, r *http.Request) {
	facebookUserID := middlewares.FBUserIDFromContext(r.Context())
	facebookUser, _ := c.userRedis.FindByID(facebookUserID)
	idsJSON := r.URL.Query().Get("ids")
	var ids []string
	// JSON unmarshal url query ids
	json.Unmarshal([]byte(idsJSON), &ids)

	returnPhotos := []*domain.WebPhoto{}
	// Find batch fb photo ids on redis.
	for _, facebookPhotoID := range ids {
		cachePhoto, err := c.photoRedis.FindByID(facebookPhotoID)
		if err == nil && (facebookUser.ID == cachePhoto.Uploader || utils.IsIn(facebookUser.ID, cachePhoto.TaggedUsers)) { // verify ACL (tagged or uploader)
			returnPhotos = append(returnPhotos, domain.WebPhotoFromCachePhoto(cachePhoto))
		}
	}

	c.render.JSON(w, http.StatusOK, returnPhotos)
}

func (c Controller) postPhotoHandler(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.FBUserIDFromContext(r.Context())
	facebookUser, _ := c.userRedis.FindByID(userID)
	webPhoto, err := FromRequest(r, facebookUser)
	if err != nil {
		c.render.JSON(w, http.StatusBadRequest, nil)
		return
	}

	cachePhoto := domain.CachePhotoFromWebPhoto(webPhoto)

	c.photoRedis.Save(cachePhoto)

	c.photoPublisher.Publish(cachePhoto)

	c.render.JSON(w, http.StatusCreated, webPhoto)
}

func (c Controller) putPhotoHandler(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.FBUserIDFromContext(r.Context())
	cacheUser, _ := c.userRedis.FindByID(userID)
	vars := mux.Vars(r)
	idPhoto := vars["id"]

	cachePhoto, err := c.photoRedis.FindByID(idPhoto)
	if err != nil {
		c.render.JSON(w, http.StatusNotFound, nil)
		return
	}

	if cachePhoto.Uploader != cacheUser.ID && !utils.IsIn(cacheUser.ID, cachePhoto.TaggedUsers) {
		c.render.JSON(w, http.StatusForbidden, nil)
		return
	}

	categories := c.categoryRedis.GetAll()
	userCategories := c.categoryRedis.GetCategoriesForUser(cacheUser)

	webPhoto, err := FromPutRequest(r, cachePhoto, categories, userCategories, cacheUser.ID)
	if err != nil {
		c.render.JSON(w, http.StatusBadRequest, nil)
		return
	}

	webCachePhoto := domain.CachePhotoFromWebPhoto(webPhoto)
	cachePhoto.Categories = webCachePhoto.Categories
	c.photoRedis.SaveCategoriesForPhoto(cachePhoto)
	c.photoRedis.SaveUserCategoriesForPhoto(cachePhoto, cacheUser)

	dbPhoto := domain.DBPhotoFromCachePhoto(cachePhoto)
	c.photoPersistPublisher.Publish(dbPhoto)

	c.render.JSON(w, http.StatusOK, webPhoto)
}
