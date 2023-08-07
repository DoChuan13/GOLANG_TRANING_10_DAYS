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

	admin, _ := sarama.NewClusterAdmin(bootstrapServers, nil)
	err = admin.CreateTopic(topic, &sarama.TopicDetail{NumPartitions: 2, ReplicationFactor: 1}, false)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		err := producer.Close()
		if err != nil {
			fmt.Println("Error closing producer: ", err)
			return
		}
		fmt.Println("AsyncProducer closed")
	}()

	parId := 0
producerLoop:
	for {
		if parId == 0 {
			parId++
			fmt.Println("Partition Id", parId)
		} else {
			parId = 0
			fmt.Println("Partition Id", parId)

		}
		value := fmt.Sprintf("Message-%d", i)
		message := sarama.ProducerMessage{Topic: topic, Partition: int32(parId), Value: sarama.StringEncoder(value)}

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
