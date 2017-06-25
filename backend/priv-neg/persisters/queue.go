package persisters

import (
	"fmt"
	"log"
	"os"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/utils"
	"github.com/streadway/amqp"
)

// Consumer - interface.
type Consumer interface {
	Consume()
}

// Publisher - interface.
type Publisher interface {
	Publish(Queueable)
	GetMessageTotal() int
}

// TotalStats - Returns totals for all of the queues.
type TotalStats struct {
	TotalMessageCount uint `json:"messageCount"`
}

// Queueable - What Publishers should accept.
type Queueable interface{}

// NewQueue - Returns a new RabbitMQ Channel and Connection.
func NewQueue(logger *log.Logger) (*amqp.Channel, *amqp.Connection) {
	if !utils.WaitForService(fmt.Sprintf("%s:%s", os.Getenv("RABBITMQ_HOSTNAME"), "5672"), logger) {
		panic("Could not find RabbitMQ..")
	}
	conn, err := amqp.Dial(
		fmt.Sprintf("amqp://%s:%s@%s:5672/",
			os.Getenv("RABBITMQ_USER"),
			os.Getenv("RABBITMQ_PASS"),
			os.Getenv("RABBITMQ_HOSTNAME"),
		),
	)
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	// defer conn.Close()
	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")

	return ch, conn
}
