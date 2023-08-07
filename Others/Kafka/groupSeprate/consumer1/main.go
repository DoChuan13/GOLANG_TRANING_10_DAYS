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
	partition, _ := strconv.Atoi(util.GetEnv(util.Partition, strconv.Itoa(1)))

	//cfg := sarama.NewConfig()

	client, err := sarama.NewClient(bootstrapServers, nil)
	om, err := sarama.NewClusterAdminFromClient(client)
	defer func() {
		_ = om.Close()
	}()

	//admin, err := sarama.NewClusterAdmin(bootstrapServers, nil)
	////admin.CreateTopic()
	//topics, err := admin.ListTopics()
	//_, exist := topics["my-topic"]
	//fmt.Println("Exist==>", exist)
	//
	//time.Sleep(2 * time.Second)

	consumer, err := sarama.NewConsumer(bootstrapServers, nil)
	if err != nil {
		panic("Error creating the consumer")
	}

	defer func() {
		err := consumer.Close()
		if err != nil {
			fmt.Println("Error closing consumer: ", err)
			return
		}
		fmt.Println("Consumer closed")
	}()

	partitionConsumer, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetOldest)

	//max := consumer.HighWaterMarks()
	////partitionConsumer, err := consumer.ConsumePartition(topic, int32(partition), 20)
	//index := partitionConsumer.HighWaterMarkOffset()

	if err != nil {
		panic("Error creating the partition consumer")
	}

	defer func() {
		err := partitionConsumer.Close()
		if err != nil {
			fmt.Println("Error closing partition consumer: ", err)
			return
		}
		fmt.Println("Consumer partition closed")
	}()

consumerLoop:
	for {
		time.Sleep(500 * time.Millisecond)
		select {
		case msg := <-partitionConsumer.Messages():
			fmt.Printf(
				"Received message from %s-%d [%d]: %s = %s\n", msg.Topic, msg.Partition, msg.Offset, string(msg.Key),
				string(msg.Value),
			)
			//fmt.Println("High Water Max (topic and partition 0) and Max Offset===>", max["my-topic"][0], index)
		case sig := <-signals:
			fmt.Println("Got signal: ", sig)
			break consumerLoop
		}
	}
}
