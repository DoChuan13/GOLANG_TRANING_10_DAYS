package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/ppatierno/kafka-go-examples/util"
)

func main() {

	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL)

	bootstrapServers := strings.Split(util.GetEnv(util.BootstrapServers, "localhost:9093"), ",")
	topic := util.GetEnv(util.Topic, "my-topic")
	delayMs, _ := strconv.Atoi(util.GetEnv(util.DelayMs, strconv.Itoa(1000)))

	producer, err := sarama.NewSyncProducer(bootstrapServers, nil)
	if err != nil {
		panic("Error creating the sync producer")
	}
	i := 1

	defer func() {
		err := producer.Close()
		if err != nil {
			fmt.Println("Error closing producer: ", err)
			return
		}
		fmt.Println("AsyncProducer closed")
	}()

producerLoop:
	for {

		value := fmt.Sprintf("Message-%d", i)
		message := sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(value)}

		partition, offset, err := producer.SendMessage(&message)
		if err != nil {
			fmt.Println("Error sending message: ", err)
		} else {
			fmt.Printf("Sent message value='%s' at partition = %d, offset = %d\n", value, partition, offset)
		}

		i++

		select {
		case sig := <-signals:
			fmt.Println("Got signal: ", sig)
			break producerLoop
		default:
			time.Sleep(time.Duration(delayMs) * time.Millisecond)
		}
	}
}
