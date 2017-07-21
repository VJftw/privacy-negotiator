package friend

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/category"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/middlewares"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

// Controller - Handles users
type Controller struct {
	logger                 *log.Logger
	render                 *render.Render
	userRedis              *user.RedisManager
	friendRedis            *RedisManager
	categoryRedis          *category.RedisManager
	friendPersistPublisher *PersistPublisher
}

// NewController - returns a new controller for users.
func NewController(
	controllerLogger *log.Logger,
	renderer *render.Render,
	userRedisManager *user.RedisManager,
	friendRedisManager *RedisManager,
	categoryRedisManager *category.RedisManager,
	friendPersistPublisher *PersistPublisher,
) *Controller {
	return &Controller{
		logger:                 controllerLogger,
		render:                 renderer,
		userRedis:              userRedisManager,
		friendRedis:            friendRedisManager,
		categoryRedis:          categoryRedisManager,
		friendPersistPublisher: friendPersistPublisher,
	}
}

// Setup - Sets up the Auth Controller
func (c Controller) Setup(router *mux.Router) {
	router.Handle("/v1/friends", negroni.New(
		middlewares.NewJWT(c.render),
		negroni.Wrap(http.HandlerFunc(c.getFriendsHandler)),
	)).Methods("GET")

	router.Handle("/v1/cliques", negroni.New(
		middlewares.NewJWT(c.render),
		negroni.Wrap(http.HandlerFunc(c.getCliquesHandler)),
	)).Methods("GET")

	router.Handle("/v1/cliques/{id}", negroni.New(
		middlewares.NewJWT(c.render),
		negroni.Wrap(http.HandlerFunc(c.putCliquesHandler)),
	)).Methods("PUT")

	log.Println("Set up Friend controller.")

}

func (c Controller) getCliquesHandler(w http.ResponseWriter, r *http.Request) {
	facebookUserID := middlewares.FBUserIDFromContext(r.Context())
	facebookUser, _ := c.userRedis.FindByID(facebookUserID)

	userCliques := c.friendRedis.GetCliquesForAUserID(facebookUser.ID)

	webCliques := []domain.WebClique{}
	for _, userClique := range userCliques {
		webCliques = append(webCliques, domain.WebCliqueFromCacheClique(userClique))
	}

	c.render.JSON(w, http.StatusOK, webCliques)
}

func (c Controller) putCliquesHandler(w http.ResponseWriter, r *http.Request) {
	facebookUserID := middlewares.FBUserIDFromContext(r.Context())
	facebookUser, _ := c.userRedis.FindByID(facebookUserID)

	vars := mux.Vars(r)
	idClique := vars["id"]

	clique, err := c.friendRedis.GetCliqueByIDAndUser(idClique, facebookUser)
	if err != nil {
		c.render.JSON(w, http.StatusNotFound, nil)
		return
	}

	categories := c.categoryRedis.GetAll()
	userCategories := c.categoryRedis.GetCategoriesForUser(facebookUser)

	clique, err = FromPutRequest(r, clique, categories, userCategories)
	if err != nil {
		c.render.JSON(w, http.StatusBadRequest, nil)
		return
	}

	c.friendRedis.AddCliqueToUserID(facebookUser.ID, clique)

	dbUserClique := domain.DBUserCliqueFromCacheCliqueAndUserID(clique, facebookUser.ID)
	c.friendPersistPublisher.Publish(dbUserClique)

	c.render.JSON(w, http.StatusOK, domain.WebCliqueFromCacheClique(*clique))
}

func (c Controller) getFriendsHandler(w http.ResponseWriter, r *http.Request) {
	facebookUserID := middlewares.FBUserIDFromContext(r.Context())
	facebookUser, _ := c.userRedis.FindByID(facebookUserID)
	idsJSON := r.URL.Query().Get("ids")
	var ids []string
	// JSON unmarshal url query ids
	json.Unmarshal([]byte(idsJSON), &ids)

	returnFriendships := []*domain.WebFriendship{}
	userCliques := c.friendRedis.GetCliqueIDsForAUserID(facebookUser.ID)
	// Find batch fb friend ids on redis.
	for _, friendUserID := range ids {
		cacheFriendship, err := c.friendRedis.FindByIDAndUser(friendUserID, facebookUser)
		if err != nil {
			// If the friendship doesn't exist, skip.
			break
		}
		// If the friendship does exist, find common cliques and return a webFriendship
		friendCliques := c.friendRedis.GetCliqueIDsForAUserID(friendUserID)
		commonCliques := []string{}
		for _, userClique := range userCliques {
			for _, friendClique := range friendCliques {
				if friendClique == userClique {
					commonCliques = append(commonCliques, userClique)
					break
				}
			}
		}

		webFriendship := domain.WebFriendshipFromCacheFriendshipAndCliques(cacheFriendship, commonCliques)
		if err == nil {
			returnFriendships = append(returnFriendships, webFriendship)
		}
	}

	c.render.JSON(w, http.StatusOK, returnFriendships)
}
