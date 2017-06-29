package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/auth"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/category"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/friend"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/photo"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/persisters"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/routers"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/routers/websocket"
	"github.com/unrolled/render"
)

// PrivNegAPI - The Privacy Negotiation API app
type PrivNegAPI struct {
	Router *routers.MuxRouter
	server *http.Server
}

// NewPrivNegAPI - Returns a new Privacy Negotiation API app
func NewPrivNegAPI() App {
	privNegAPI := &PrivNegAPI{}

	wsLogger := log.New(os.Stdout, "[websocket] ", log.Lshortfile)
	queueLogger := log.New(os.Stdout, "[queue] ", log.Lshortfile)
	cacheLogger := log.New(os.Stdout, "[cache] ", log.Lshortfile)
	controllerLogger := log.New(os.Stdout, "[controller]", log.Lshortfile)

	redisCache := persisters.NewRedisDB(cacheLogger)
	rabbitMQ, _ := persisters.NewQueue(queueLogger)

	userRedisManager := user.NewRedisManager(cacheLogger, redisCache)
	photoRedisManager := photo.NewRedisManager(cacheLogger, redisCache)
	friendRedisManager := friend.NewRedisManager(cacheLogger, redisCache)
	categoryRedisManager := category.NewRedisManager(cacheLogger, redisCache)

	authPublisher := auth.NewPublisher(queueLogger, rabbitMQ)
	syncPublisher := photo.NewPublisher(queueLogger, rabbitMQ)
	categoryPublisher := category.NewPublisher(queueLogger, rabbitMQ)
	friendPublisher := friend.NewPublisher(queueLogger, rabbitMQ)
	friendCliquePersistPublisher := friend.NewPersistPublisher(queueLogger, rabbitMQ)

	renderer := render.New()

	authController := auth.NewController(controllerLogger, renderer, authPublisher, userRedisManager)
	userController := user.NewController(controllerLogger, renderer, userRedisManager)
	photoController := photo.NewController(controllerLogger, renderer, photoRedisManager, userRedisManager, categoryRedisManager, syncPublisher)
	categoryController := category.NewController(controllerLogger, renderer, userRedisManager, categoryRedisManager, categoryPublisher)
	friendController := friend.NewController(controllerLogger, renderer, userRedisManager, friendRedisManager, categoryRedisManager, friendCliquePersistPublisher)
	websocketController := websocket.NewController(wsLogger, renderer, redisCache)
	healthController := routers.NewHealthController(controllerLogger, renderer, []persisters.Publisher{
		authPublisher,
		syncPublisher,
		categoryPublisher,
		friendPublisher,
	})

	privNegAPI.Router = routers.NewMuxRouter([]routers.Routable{
		authController,
		userController,
		photoController,
		categoryController,
		friendController,
		websocketController,
		healthController,
	}, true)

	port := os.Getenv("PORT")
	privNegAPI.server = &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		Handler:        privNegAPI.Router.Negroni,
		ReadTimeout:    1 * time.Hour,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return privNegAPI
}

// Stop - Stops the API
func (p *PrivNegAPI) Stop() {
	if err := p.server.Shutdown(nil); err != nil {
		panic(err)
	}
}

// Start - Starts the API
func (p *PrivNegAPI) Start() {
	if err := p.server.ListenAndServe(); err != nil {
		log.Printf("Error %s", err)
	}
}
