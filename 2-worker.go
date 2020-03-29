package main

import (
	"bytes"
	"log"
	"time"

	"github.com/google/uuid"
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

	// make sure one receive only receive one message at a time. Will get new one once current one is done
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	utils.FailOnError(err, "Failed to set QoS")

	q, err := ch.QueueDeclare(
		utils.WorkQueueName, // name
		true,                // durable. Set to true in order to resist data lose if rabbitmq restarted
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")

	consumerName := uuid.New().String()
	msgs, err := ch.Consume(
		q.Name,       // queue
		consumerName, // consumer
		// auto-ack. Set to false in order to make sure server will delete message only it receives manualy ack. Resist data lost if receiver restarts
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	utils.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("%s Received a message: %s", consumerName, d.Body)
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Printf("Done")
			// Ack the message has been processed, send back to rabbitmq server. Resist data lost if receiver restarts
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	<-forever
}
