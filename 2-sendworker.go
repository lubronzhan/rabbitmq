package main

import (
	"os"

	"github.com/lubronzhan/rabbitmp/pkg/utils"
	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		utils.WorkQueueName, // name
		true,                // durable. Set to true in order to resist data lost if rabbitmq restarts
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")

	message := os.Args[1]
	err = ch.Publish(
		"", // exchange. Empty, it will use the default exchange
		q.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			// Resist data lost if rabbitmq restarts
			ContentType: "text/plain",
			Body:        []byte(message),
		})

	utils.FailOnError(err, "Failed to publish a message")
}
