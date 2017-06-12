package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/auth"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/category"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/photo"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/persisters"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/routers"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/routers/websocket"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/utils"
	"github.com/unrolled/render"
)

// PrivNegAPI - The Privacy Negotiation API app
type PrivNegAPI struct {
	Router *routers.MuxRouter
}

// NewPrivNegAPI - Returns a new Privacy Negotiation API app
func NewPrivNegAPI() {
	privNegAPI := &PrivNegAPI{}

	wsLogger := log.New(os.Stdout, "[websocket] ", log.Lshortfile)
	queueLogger := log.New(os.Stdout, "[queue] ", log.Lshortfile)
	cacheLogger := log.New(os.Stdout, "[cache] ", log.Lshortfile)
	controllerLogger := log.New(os.Stdout, "[controller]", log.Lshortfile)

	redisCache := persisters.NewRedisDB(cacheLogger)

	userManager := user.NewAPIManager(cacheLogger, redisCache)
	photoManager := photo.NewAPIManager(cacheLogger, redisCache)
	categoryManager := category.NewAPIManager(cacheLogger, redisCache)

	authQueue := auth.NewAuthQueue(queueLogger, userManager)
	syncQueue := photo.NewSyncQueue(queueLogger, photoManager, userManager)

	renderer := render.New()

	authController := auth.NewController(controllerLogger, renderer, authQueue, userManager)
	userController := user.NewController(controllerLogger, renderer, userManager)
	photoController := photo.NewController(controllerLogger, renderer, photoManager, syncQueue)
	categoryController := category.NewController(controllerLogger, renderer, categoryManager)
	websocketController := websocket.NewController(wsLogger, renderer)

	privNegAPI.Router = routers.NewMuxRouter([]routers.Routable{
		authController,
		userController,
		photoController,
		categoryController,
		websocketController,
	}, true)

	// Initialise queues
	utils.SetupQueues([]utils.DeclarableQueue{
		authQueue,
		syncQueue,
	}, queueLogger)

	port := os.Getenv("PORT")
	s := &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		Handler:        privNegAPI.Router.Negroni,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
