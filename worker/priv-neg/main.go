package main

import (
	"fmt"
	"log"
	"os"

	"github.com/VJftw/privacy-negotiator/worker/priv-neg/domain/user"
	"github.com/VJftw/privacy-negotiator/worker/priv-neg/persisters"
	"github.com/VJftw/privacy-negotiator/worker/priv-neg/queues"
	"github.com/facebookgo/inject"
)

// PrivNegWorker - The Privacy Negotiation API app
type PrivNegWorker struct {
	Graph *inject.Graph
}

// NewPrivNegWorker - Returns a new Privacy Negotiation API app
func NewPrivNegWorker() *PrivNegWorker {
	privNegWorker := PrivNegWorker{
		Graph: &inject.Graph{},
	}

	mainLogger := log.New(os.Stdout, "[main] ", log.Lshortfile)
	dbLogger := log.New(os.Stdout, "[database] ", log.Lshortfile)
	queueLogger := log.New(os.Stdout, "[queue] ", log.Lshortfile)

	// Initialise persisters to pass into managers
	gormDB := persisters.NewGORMDB(
		dbLogger,
		&user.FacebookUser{},
	)

	qGetFacebookLongLivedToken := queues.NewGetFacebookLongLivedToken()

	err := privNegWorker.Graph.Provide(
		&inject.Object{Name: "logger.main", Value: mainLogger},
		&inject.Object{Name: "logger.db", Value: dbLogger},
		&inject.Object{Name: "logger.queue", Value: queueLogger},
		&inject.Object{Name: "user.manager", Value: user.NewManager(gormDB)},
		&inject.Object{Name: "queues.getFacebookLongLivedToken", Value: qGetFacebookLongLivedToken},
	)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := privNegWorker.Graph.Populate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Initialise queues
	queues.StartQueue(qGetFacebookLongLivedToken, queueLogger)

	return &privNegWorker
}

func main() {
	NewPrivNegWorker()
}
