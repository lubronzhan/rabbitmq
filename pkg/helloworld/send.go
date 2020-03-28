package helloworld

import (
	"github.com/streadway/amqp"
)

func Connect() {
	conn, err := amqp.Dail("amqp://guest:guest@localhost:5672/")
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
}
