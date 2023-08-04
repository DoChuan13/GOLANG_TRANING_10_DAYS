package kafka

import (
	"Mock_Project/model"
	"fmt"
	"github.com/IBM/sarama"
	"os"
	"strings"
	"sync"
)

type kafka struct {
	kafkaClient map[string]model.Kafka
	config      *sarama.Config
	admin       sarama.ClusterAdmin
	system      *model.Server
}

// NewKafkaHandler constructor
func NewKafkaHandler(cfg *model.Server) (IKafkaHandler, error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Producer.MaxMessageBytes = 10e6
	kafkaConfig.Producer.Retry.Max = 5
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	admin, err := sarama.NewClusterAdmin(cfg.Broker, kafkaConfig)
	if err != nil {
		return nil, err
	}
	return &kafka{
		kafkaClient: make(map[string]model.Kafka),
		config:      kafkaConfig,
		system:      cfg,
		admin:       admin,
	}, nil
}

func (k kafka) InitConnection(topic string) error {
	_, isExists := k.kafkaClient[topic]
	if isExists {
		return nil
	}
	err := k.CreateTopic(topic, 1)
	if err != nil {
		return err
	}
	consumer, err := sarama.NewConsumer(k.system.Broker, k.config)
	if err != nil {
		return err
	}

	producer, err := sarama.NewSyncProducer(k.system.Broker, k.config)
	if err != nil {
		return err
	}
	newKafka := model.Kafka{Consumer: consumer, Producer: producer}
	k.kafkaClient[topic] = newKafka
	return nil
}

func (k kafka) CloseTopic() error {
	for _, client := range k.kafkaClient {
		if err := client.Consumer.Close(); err != nil {
			return err
		}
		if err := client.Producer.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (k kafka) CreateTopic(topic string, partitionNum int32) error {
	// Set number of partitions for topic
	err := k.admin.CreateTopic(
		topic, &sarama.TopicDetail{
			NumPartitions:     partitionNum,
			ReplicationFactor: 1,
		}, false,
	)

	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			_ = k.admin.DeleteTopic(topic)
			return k.CreateTopic(topic, partitionNum)
		}
		fmt.Println("Error setting number of partitions:=> ", err)
		os.Exit(1)
	}
	return nil
}

func (k kafka) ProducerData(topic string, partition int32, content string) error {
	currentClient, isExists := k.kafkaClient[topic]
	if !isExists {
		return fmt.Errorf("kafka topic not found")
	}

	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: partition,
		Value:     sarama.StringEncoder(content),
	}

	partition, _, err := currentClient.Producer.SendMessage(msg)
	if err != nil {
		return err
	}

	/*fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)*/
	return nil
}

func (k kafka) ConsumerData(topic string, partition int32, parse ParseStruct) (
	[]string, error,
) {
	currentClient, isExists := k.kafkaClient[topic]
	if !isExists {
		return nil, fmt.Errorf("kafka topic not found")
	}

	resp, err := currentClient.Consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
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

func (k kafka) RemoveTopic() error {
	defer func() { _ = k.admin.Close() }()
	for topic := range k.system.Topics {
		err := k.admin.DeleteTopic(topic)
		if err != nil {
			return err
		}
	}
	return nil
}
