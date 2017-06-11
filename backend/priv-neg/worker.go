package main

import (
	"log"
	"os"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/auth"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/category"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/photo"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/persisters"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/utils"
)

// NewPrivNegWorker - Returns a new Privacy Negotiation API app
func NewPrivNegWorker(queue string) {

	dbLogger := log.New(os.Stdout, "[database] ", log.Lshortfile)
	queueLogger := log.New(os.Stdout, "[queue] ", log.Lshortfile)
	cacheLogger := log.New(os.Stdout, "[cache] ", log.Lshortfile)

	// Initialise persisters to pass into managers
	gormDB := persisters.NewGORMDB(
		dbLogger,
		&user.FacebookUser{},
		&photo.FacebookPhoto{},
		&category.Category{},
	)
	redisCache := persisters.NewRedisDB(cacheLogger)

	userManager := user.NewWorkerManager(dbLogger, gormDB, cacheLogger, redisCache)
	photoManager := photo.NewWorkerManager(dbLogger, gormDB, cacheLogger, redisCache)

	var q utils.DeclarableQueue
	switch queue {
	case "auth-queue":
		q = auth.NewAuthQueue(queueLogger, userManager)
		break
	case "sync-queue":
		q = photo.NewSyncQueue(queueLogger, photoManager, userManager)
		break
	default:
		panic("Invalid queue selected")
	}

	utils.SetupQueues([]utils.DeclarableQueue{q}, queueLogger)
	utils.Consume([]utils.DeclarableQueue{q}, queueLogger)

}
