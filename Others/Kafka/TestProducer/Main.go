package main

import (
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"os"
	"sync"
	"time"
)

type Kafka struct {
	asyncProducer sarama.AsyncProducer
	syncProducer  sarama.SyncProducer
	consumer      sarama.Consumer
}

func Consumer(consumer sarama.Consumer, topic string) error {
	signals := make(chan os.Signal, 1)
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
consumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			fmt.Printf(
				"Received message from %s-%d [%d]: %s = %s\n", msg.Topic, msg.Partition, msg.Offset, string(msg.Key),
				string(msg.Value),
			)
		case sig := <-signals:
			fmt.Println("Got signal: ", sig)
			break consumerLoop
		}
	}
	return nil
}

func AsyncProducer(asyncProducer sarama.AsyncProducer, topic string, content string) error {
	message := &sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(content)}
	asyncProducer.Input() <- message
	var err error = nil
loop:
	for {
		select {
		case success := <-asyncProducer.Successes():
			fmt.Println("Success", success)
			err = nil
			break loop
		case errInf := <-asyncProducer.Errors():
			err = errInf
			break loop
		}
	}
	fmt.Println("Err info==>", err)
	return err
}

func SyncProducer(syncProducer sarama.SyncProducer, topic string, content string) error {
	message := &sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(content)}
	_, _, err := syncProducer.SendMessage(message)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	var wg = sync.WaitGroup{}

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producerAsync, err := sarama.NewAsyncProducer([]string{"localhost:9093"}, config)
	if err != nil {
		panic(err)
	}

	producerSync, err := sarama.NewSyncProducer([]string{"localhost:9093"}, config)
	if err != nil {
		panic(err)
	}
	consumer, err := sarama.NewConsumer([]string{"localhost:9093"}, config)

	if err != nil {
		panic(err)
	}
	defer func() {
		if err := producerAsync.Close(); err != nil {
			log.Fatalln(err)
		}
		if err := producerSync.Close(); err != nil {
			log.Fatalln(err)
		}
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	for i := 0; i < 10; i++ {
		wg.Add(2)
		go processSync(producerSync, i, &wg)
		go processAsync(producerAsync, i, &wg)
		//go processSyncInner(i, &wg)
	}

	wg.Wait()
	fmt.Println("Done Async")
	fmt.Println("Done Sync")

	fmt.Println("Sleep 1s")
	time.Sleep(time.Second)

	signals := make(chan os.Signal, 1)
	//wg.Add(2)
	consumerAsyncTopic(&wg, err, consumer, signals)
	go consumerSyncTopic(&wg, err, consumer, signals)
	fmt.Println("Waiting")
	wg.Wait()
	fmt.Println("Finished")
}

func consumerAsyncTopic(wg *sync.WaitGroup, err error, consumer sarama.Consumer, signals chan os.Signal) {
	//defer wg.Done()
	topic := "Hello-world-Async"
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
consumerLoop:
	for {
		fmt.Println("In Async Consumer")
		select {
		case msg := <-partitionConsumer.Messages():
			fmt.Printf(
				"Received message from %s-%d [%d]: %s = %s\n", msg.Topic, msg.Partition, msg.Offset, string(msg.Key),
				string(msg.Value),
			)
			//break consumerLoop
		case sig := <-signals:
			fmt.Println("Got signal: ", sig)
			break consumerLoop
		}
	}
}

func consumerSyncTopic(wg *sync.WaitGroup, err error, consumer sarama.Consumer, signals chan os.Signal) {
	defer wg.Done()
	topic := "Hello-world-Sync"
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
consumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			fmt.Printf(
				"Received message from %s-%d [%d]: %s = %s\n", msg.Topic, msg.Partition, msg.Offset, string(msg.Key),
				string(msg.Value),
			)
			break consumerLoop
		case sig := <-signals:
			fmt.Println("Got signal: ", sig)
			break consumerLoop
		}
	}
}

func processAsync(producer sarama.AsyncProducer, index int, wg *sync.WaitGroup) {
	time.Sleep(2 * time.Second)
	defer wg.Done()
	topic := "Hello-world-Async"
	var msg = "Hello World Async! ==>" + string(rune(index))
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	}
loop:
	for {
		producer.Input() <- message

		select {
		case success := <-producer.Successes():
			fmt.Println("Message sent Async:", success.Offset)
			break loop
		case err := <-producer.Errors():
			fmt.Println("Failed to send message Async:", err)
			break loop
		}
	}
}

func processSync(producer sarama.SyncProducer, index int, wg *sync.WaitGroup) {
	time.Sleep(2 * time.Second)
	defer wg.Done()
	topic := "Hello-world-Sync"
	var msg = "Hello World Sync! ==>" + string(rune(index))
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	}

	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		fmt.Println("Failed to send message Sync:", err)
		return
	}
	fmt.Println("Message sent Sync (partition, offset):", partition, offset)
}

func processSyncInner(index int, wg *sync.WaitGroup) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	producer, err := sarama.NewSyncProducer([]string{"localhost:9093"}, config)
	if err != nil {
		panic(err)
	}

	defer wg.Done()
	topic := "Hello-world-Sync"
	var msg = "Hello World Sync! ==>" + string(rune(index))
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	}

	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		fmt.Println("Failed to send message Sync:", err)
		return
	}
	fmt.Println("Message sent Sync (partition, offset):", partition, offset)
}
