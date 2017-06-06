package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/VJftw/privacy-negotiator/api/priv-neg/domain/auth"
	"github.com/VJftw/privacy-negotiator/api/priv-neg/queues"
	"github.com/VJftw/privacy-negotiator/api/priv-neg/routers"
	"github.com/facebookgo/inject"
)

type PrivNegApi struct {
	Graph  *inject.Graph
	Router *routers.MuxRouter
}

func NewPrivNegApi() *PrivNegApi {
	PrivNegApi := PrivNegApi{
		Graph: &inject.Graph{},
	}

	mainLogger := log.New(os.Stdout, "[main] ", log.Lshortfile)
	wsLogger := log.New(os.Stdout, "[websocket] ", log.Lshortfile)
	dbLogger := log.New(os.Stdout, "[database] ", log.Lshortfile)
	queueLogger := log.New(os.Stdout, "[queue] ", log.Lshortfile)

	var authController auth.Controller
	qGetFacebookLongLivedToken := queues.NewGetFacebookLongLivedToken()

	err := PrivNegApi.Graph.Provide(
		&inject.Object{Name: "logger.main", Value: mainLogger},
		&inject.Object{Name: "logger.ws", Value: wsLogger},
		&inject.Object{Name: "logger.db", Value: dbLogger},
		&inject.Object{Name: "auth.resolver", Value: auth.NewResolver()},
		&inject.Object{Name: "auth.provider", Value: auth.NewProvider()},
		&inject.Object{Name: "auth.graphAPI", Value: auth.NewGraphAPI()},
		&inject.Object{Name: "auth.controller", Value: &authController},
		&inject.Object{Name: "queues.getFacebookLongLivedToken", Value: qGetFacebookLongLivedToken},
	)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := PrivNegApi.Graph.Populate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	muxRouter := routers.NewMuxRouter([]routers.Routable{
		&authController,
	}, true)

	PrivNegApi.Router = muxRouter

	// Initialise queues
	queues.SetupQueues([]queues.DeclarableQueue{
		qGetFacebookLongLivedToken,
	}, queueLogger)

	return &PrivNegApi
}

func main() {
	app := NewPrivNegApi()
	port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), app.Router.Handler))
}
