package main

import (
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {
	c, err := kafka.NewConsumer(
		&kafka.ConfigMap{
			"bootstrap.servers": "localhost:9093",
			"group.id":          "myGroup",
			"auto.offset.reset": "earliest",
			//"session.timeout.ms": "130",
		},
	)

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{"myTopic"}, nil)

	// A signal handler or similar could be used to set this to false to break the loop.
	run := true

	for run {
		time.Sleep(time.Second)
		msg, err := c.ReadMessage(time.Nanosecond)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else if !err.(kafka.Error).IsTimeout() {
			// The client will automatically try to recover from all errors.
			// Timeout is not considered an error because it is raised by
			// ReadMessage in absence of messages.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		} else {
			fmt.Println("Time Out ==>", err.Error())
		}
	}

	c.Close()
}
