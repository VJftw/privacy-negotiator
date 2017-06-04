package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/VJftw/privacy-negotiator/api/priv-neg/domain/auth"
	"github.com/VJftw/privacy-negotiator/api/priv-neg/routers"
	"github.com/facebookgo/inject"
)

type PrivNegApp struct {
	Graph  *inject.Graph
	Router *routers.MuxRouter
}

func NewPrivNegApp() *PrivNegApp {
	privNegApp := PrivNegApp{
		Graph: &inject.Graph{},
	}

	mainLogger := log.New(os.Stdout, "[main] ", log.Lshortfile)
	wsLogger := log.New(os.Stdout, "[websocket] ", log.Lshortfile)
	dbLogger := log.New(os.Stdout, "[database] ", log.Lshortfile)

	var authController auth.Controller

	err := privNegApp.Graph.Provide(
		&inject.Object{Name: "logger.main", Value: mainLogger},
		&inject.Object{Name: "logger.ws", Value: wsLogger},
		&inject.Object{Name: "logger.db", Value: dbLogger},
		&inject.Object{Name: "auth.resolver", Value: auth.NewResolver()},
		&inject.Object{Name: "auth.provider", Value: auth.NewProvider()},
		&inject.Object{Name: "auth.graphAPI", Value: auth.NewGraphAPI()},
		&inject.Object{Name: "auth.controller", Value: &authController},
	)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := privNegApp.Graph.Populate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	muxRouter := routers.NewMuxRouter([]routers.Routable{
		&authController,
	}, true)

	privNegApp.Router = muxRouter

	return &privNegApp
}

func main() {
	app := NewPrivNegApp()
	port := os.Getenv("HTTP_PORT")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), app.Router.Handler))
}
