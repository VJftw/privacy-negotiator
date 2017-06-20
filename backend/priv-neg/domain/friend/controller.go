package friend

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/middlewares"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

// Controller - Handles users
type Controller struct {
	logger      *log.Logger
	render      *render.Render
	userRedis   *user.RedisManager
	friendRedis *RedisManager
}

// NewController - returns a new controller for users.
func NewController(
	controllerLogger *log.Logger,
	renderer *render.Render,
	userRedisManager *user.RedisManager,
	friendRedisManager *RedisManager,
) *Controller {
	return &Controller{
		logger:      controllerLogger,
		render:      renderer,
		userRedis:   userRedisManager,
		friendRedis: friendRedisManager,
	}
}

// Setup - Sets up the Auth Controller
func (c Controller) Setup(router *mux.Router) {
	router.Handle("/v1/friends", negroni.New(
		middlewares.NewJWT(c.render),
		negroni.Wrap(http.HandlerFunc(c.getFriendsHandler)),
	)).Methods("GET")

	router.Handle("/v1/friends", negroni.New(
		middlewares.NewJWT(c.render),
		negroni.Wrap(http.HandlerFunc(c.postFriendsHandler)),
	)).Methods("POST")

	log.Println("Set up Friend controller.")

}

func (c Controller) getFriendsHandler(w http.ResponseWriter, r *http.Request) {
	facebookUserID := middlewares.FBUserIDFromContext(r.Context())
	facebookUser, _ := c.userRedis.FindByID(facebookUserID)
	idsJSON := r.URL.Query().Get("ids")
	var ids []string
	// JSON unmarshal url query ids
	json.Unmarshal([]byte(idsJSON), &ids)

	returnFriendships := []*domain.WebFriendship{}
	// Find batch fb photo ids on redis.
	for _, friendUserID := range ids {
		cacheFriendship, err := c.friendRedis.FindByIDAndUser(friendUserID, facebookUser)
		if err != nil {
			break
		}
		webFriendship := domain.WebFriendshipFromCacheFriendshipAndCacheUser(cacheFriendship, facebookUser)
		if err == nil {
			returnFriendships = append(returnFriendships, webFriendship)
		}
	}

	c.render.JSON(w, http.StatusOK, returnFriendships)

}

func (c Controller) postFriendsHandler(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.FBUserIDFromContext(r.Context())
	cacheUser, _ := c.userRedis.FindByID(userID)
	webFriendship, err := FromRequest(r, cacheUser)
	if err != nil {
		c.render.JSON(w, http.StatusBadRequest, nil)
		return
	}

	cacheFriendship := domain.CacheFriendshipFromWebFriendship(webFriendship)

	c.friendRedis.Save(cacheUser, cacheFriendship)

	// This Queue can determine clique and tie-strength
	//c.friendPublisher.Publish(webFriendship)

	c.render.JSON(w, http.StatusCreated, webFriendship)

}
