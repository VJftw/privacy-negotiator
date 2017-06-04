package queues

import "github.com/streadway/amqp"

type DeclarableQueue interface {
	Setup(*amqp.Channel)
	Publish(Queueable)
	Consume()
}

type Queueable interface {
}
