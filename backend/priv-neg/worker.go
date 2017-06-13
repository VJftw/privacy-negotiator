package main

import (
	"log"
	"os"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/auth"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/category"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/photo"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/user"
	"github.com/VJftw/privacy-negotiator/backend/priv-neg/persisters"
	"github.com/streadway/amqp"
)

// PrivNegWorker - Holds the Worker
type PrivNegWorker struct {
	queue   persisters.Consumer
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewPrivNegWorker - Returns a new Privacy Negotiation API app
func NewPrivNegWorker(queue string) App {

	privNegWorker := &PrivNegWorker{}

	dbLogger := log.New(os.Stdout, "[database] ", log.Lshortfile)
	queueLogger := log.New(os.Stdout, "[queue] ", log.Lshortfile)
	cacheLogger := log.New(os.Stdout, "[cache] ", log.Lshortfile)

	// Initialise persisters to pass into managers
	gormDB := persisters.NewGORMDB(
		dbLogger,
		&domain.DBUser{},
		&domain.DBPhoto{},
		&domain.DBCategory{},
	)
	redisCache := persisters.NewRedisDB(cacheLogger)

	photoRedisManager := photo.NewRedisManager(cacheLogger, redisCache)
	userDBManager := user.NewDBManager(dbLogger, gormDB)
	userRedisManager := user.NewRedisManager(cacheLogger, redisCache)
	categoryDBManager := category.NewDBManager(dbLogger, gormDB)

	rabbitMQ, conn := persisters.NewQueue(queueLogger)
	var q persisters.Consumer
	switch queue {
	case "auth-long-token":
		q = auth.NewConsumer(queueLogger, rabbitMQ, userDBManager)
		break
	case "photo-tags":
		q = photo.NewConsumer(queueLogger, rabbitMQ, userDBManager, userRedisManager, photoRedisManager)
		break
	case "category-persist":
		q = category.NewConsumer(queueLogger, rabbitMQ, categoryDBManager)
		break
	default:
		panic("Invalid queue selected")
	}

	privNegWorker.channel = rabbitMQ
	privNegWorker.conn = conn
	privNegWorker.queue = q

	return privNegWorker
}

// Start - Starts the Worker
func (p *PrivNegWorker) Start() {
	p.queue.Consume()
}

// Stop - Stops the worker
func (p *PrivNegWorker) Stop() {
	p.channel.Close()
	p.conn.Close()
}
