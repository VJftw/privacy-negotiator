package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/auth"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/persisters"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/queues"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/routers"
	"github.com/facebookgo/inject"
)

// PrivNegAPI - The Privacy Negotiation API app
type PrivNegAPI struct {
	Graph  *inject.Graph
	Router *routers.MuxRouter
}

// NewPrivNegAPI - Returns a new Privacy Negotiation API app
func NewPrivNegAPI() {
	privNegAPI := PrivNegAPI{
		Graph: &inject.Graph{},
	}

	mainLogger := log.New(os.Stdout, "[main] ", log.Lshortfile)
	wsLogger := log.New(os.Stdout, "[websocket] ", log.Lshortfile)
	dbLogger := log.New(os.Stdout, "[database] ", log.Lshortfile)
	queueLogger := log.New(os.Stdout, "[queue] ", log.Lshortfile)
	cacheLogger := log.New(os.Stdout, "[cache] ", log.Lshortfile)

	var authController auth.Controller
	var userController user.Controller
	qGetFacebookLongLivedToken := queues.NewGetFacebookLongLivedToken()

	redisCache := persisters.NewRedisDB(cacheLogger)

	err := privNegAPI.Graph.Provide(
		&inject.Object{Name: "logger.main", Value: mainLogger},
		&inject.Object{Name: "logger.ws", Value: wsLogger},
		&inject.Object{Name: "logger.db", Value: dbLogger},
		&inject.Object{Name: "logger.cache", Value: cacheLogger},
		&inject.Object{Name: "logger.queue", Value: queueLogger},
		&inject.Object{Name: "auth.resolver", Value: auth.NewResolver()},
		&inject.Object{Name: "auth.provider", Value: auth.NewProvider()},
		&inject.Object{Name: "auth.graphAPI", Value: auth.NewGraphAPI()},
		&inject.Object{Name: "auth.controller", Value: &authController},
		&inject.Object{Name: "user.controller", Value: &userController},
		&inject.Object{Name: "user.manager", Value: user.NewAPIManager()},
		&inject.Object{Name: "persister.cache", Value: redisCache},
		&inject.Object{Name: "queues.getFacebookLongLivedToken", Value: qGetFacebookLongLivedToken},
	)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := privNegAPI.Graph.Populate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	muxRouter := routers.NewMuxRouter([]routers.Routable{
		&authController,
		&userController,
	}, true)

	privNegAPI.Router = muxRouter

	// Initialise queues
	queues.SetupQueues([]queues.DeclarableQueue{
		qGetFacebookLongLivedToken,
	}, queueLogger)

	port := os.Getenv("PORT")
	// privNegAPI.Router.Handler.Run(fmt.Sprintf(":%s", port))

	s := &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		Handler:        privNegAPI.Router.Negroni,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
