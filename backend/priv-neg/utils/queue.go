package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

// DeclarableQueue - Implement this to add Queues.
type DeclarableQueue interface {
	Setup(*amqp.Channel)
	Publish(Queueable)
	Consume()
}

// Queueable - Implement this to have queueable structs.
type Queueable interface {
}

// SetupQueues - Sets up the given Queue.
func SetupQueues(queues []DeclarableQueue, logger *log.Logger) (*amqp.Connection, *amqp.Channel) {
	if !WaitForService(fmt.Sprintf("%s:%s", os.Getenv("RABBITMQ_HOSTNAME"), "5672"), logger) {
		panic("Could not find RabbitMQ..")
	}
	conn, err := amqp.Dial(
		fmt.Sprintf("amqp://%s:%s@%s:5672/",
			os.Getenv("RABBITMQ_USER"),
			os.Getenv("RABBITMQ_PASS"),
			os.Getenv("RABBITMQ_HOSTNAME"),
		),
	)
	FailOnError(err, "Failed to connect to RabbitMQ")
	// defer conn.Close()
	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	// defer ch.Close()

	for _, queue := range queues {
		queue.Setup(ch)
	}

	return conn, ch
}

// Consume - Starts a given Queue Consumer
func Consume(queues []DeclarableQueue, logger *log.Logger) {
	// logger.Printf("Starting consumer for %s", queue.GetName())
	for _, queue := range queues {
		queue.Consume()
	}
}
