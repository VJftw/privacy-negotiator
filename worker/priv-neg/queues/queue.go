package queues

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/streadway/amqp"
)

type DeclarableQueue interface {
	Setup(*amqp.Channel)
	Publish(Queueable)
	Consume()
}

type Queueable interface {
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

// StartQueue - Sets up the given Queue.
func StartQueue(queue DeclarableQueue, logger *log.Logger) {
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
	failOnError(err, "Failed to connect to RabbitMQ")
	// defer conn.Close()
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	// defer ch.Close()

	queue.Setup(ch)
}

func WaitForService(address string, logger *log.Logger) bool {

	for i := 0; i < 12; i++ {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			logger.Println("Connection error:", err)
		} else {
			conn.Close()
			logger.Println(fmt.Sprintf("Connected to %s", address))
			return true
		}
		time.Sleep(5 * time.Second)
	}

	return false
}
