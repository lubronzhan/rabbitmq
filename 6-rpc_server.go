package main

import (
	"log"
	"strconv"

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
		"rpc_queue", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	utils.FailOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	utils.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			n, err := strconv.Atoi(string(d.Body))
			utils.FailOnError(err, "Failed to convert body to integer")

			log.Printf(" [.] fib(%d)", n)
			response := utils.Fib(n)

			err = ch.Publish(
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId, // map to the request id
					Body:          []byte(strconv.Itoa(response)),
				})
			utils.FailOnError(err, "Failed to publish a message")

			d.Ack(false)
		}
	}()

	log.Printf(" [*] Awaiting RPC requests")
	<-forever
}
