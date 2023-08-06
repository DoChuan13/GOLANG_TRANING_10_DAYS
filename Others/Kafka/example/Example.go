package main

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

func main() {
	topic := "my-topic"

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer([]string{"localhost:9093"}, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder("Hello World!"),
	}

	for i := 0; i < 10; i++ {
	loop:
		for {
			time.Sleep(1 * time.Second)

			producer.Input() <- message

			select {
			case success := <-producer.Successes():
				fmt.Println("Message sent:", success.Offset)
				break loop
			case err := <-producer.Errors():
				fmt.Println("Failed to send message:", err)
				break loop
			}
		}
	}
}
