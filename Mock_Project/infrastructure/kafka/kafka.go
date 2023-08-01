package kafka

import (
	"Mock_Project/model"
	"fmt"
	"github.com/IBM/sarama"
	"sync"
)

type kafka struct {
	config *sarama.Config
	system *model.KafkaSystem
}

// NewKafkaHandler constructor
func NewKafkaHandler(cfg *model.KafkaSystem) (IKafkaHandler, error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Producer.Retry.Max = 5
	kafkaConfig.Producer.Return.Successes = true
	return &kafka{
		config: kafkaConfig,
		system: cfg,
	}, nil
}

func (k kafka) ProducerData(broker []string, topic string, partition int32, content string) error {
	producer, err := sarama.NewSyncProducer(broker, k.config)
	if err != nil {
		return err
	}

	defer func() {
		_ = producer.Close()
	}()

	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: partition,
		Value:     sarama.StringEncoder(content),
	}

	partition, _, err = producer.SendMessage(msg)
	if err != nil {
		return err
	}

	/*fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)*/
	return nil
}

func (k kafka) ConsumerData(broker []string, topic string, partition int32, parse ParseStruct) (
	[]interface{}, error,
) {
	consumer, err := sarama.NewConsumer(broker, k.config)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = consumer.Close()
	}()

	resp, err := consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
	if err != nil {
		return nil, err
	}
	size := resp.HighWaterMarkOffset()
	collection := make(chan sarama.ConsumerMessage, size)
	isContinue := true
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for isContinue {
			select {
			case err := <-resp.Errors():
				fmt.Println(err)
			case msg := <-resp.Messages():
				collection <- *msg
				if len(collection) == cap(collection) {
					isContinue = false
				}
			}
		}
		close(collection)
		wg.Done()
	}()
	wg.Wait()
	result, err := parse(&collection)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (k kafka) ClearTopic(server *model.Server) error {
	config := sarama.NewConfig()
	admin, err := sarama.NewClusterAdmin(server.Broker, config)
	if err != nil {
		fmt.Println("Clear Topic Error")
		return err
	}
	defer func() { _ = admin.Close() }()

	for _, topic := range server.Topics {
		err := admin.DeleteTopic(topic)
		if err != nil {
			return err
		}
	}

	fmt.Println("Topic deleted successfully!")
	return nil
}
