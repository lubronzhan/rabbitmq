package utils

import (
	"log"
)

const (
	HelloWorldQueueName = "hello"
	WorkQueueName       = "workqueue"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err.Error())
	}
}
