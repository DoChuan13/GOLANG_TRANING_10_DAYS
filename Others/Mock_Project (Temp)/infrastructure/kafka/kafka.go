package kafka

import (
	"Mock_Project/model"
	"fmt"
	"github.com/IBM/sarama"
	"os"
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
	config := sarama.NewConfig()

	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.MaxMessageBytes = 10e6
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	var err error = nil
	admin, err := sarama.NewClusterAdmin(cfg.Broker, config)
	if err != nil {
		return nil, err
	}
	return &kafka{
		kafkaClient: make(map[string]model.Kafka),
		config:      config,
		system:      cfg,
		admin:       admin,
	}, err
}

func (k kafka) InitConnection(topic string) error {
	var err error = nil
	_, isExists := k.kafkaClient[topic]
	if isExists {
		return err
	}
	//err = k.CreateTopic(topic, 1)
	//if err != nil {
	//	return err
	//}
	syncProducer, err := sarama.NewSyncProducer(k.system.Broker, k.config)
	if err != nil {
		return err
	}
	aSyncProducer, err := sarama.NewAsyncProducer(k.system.Broker, k.config)
	if err != nil {
		return err
	}
	consumer, err := sarama.NewConsumer(k.system.Broker, k.config)
	if err != nil {
		return err
	}
	newKafka := model.Kafka{
		Consumer:      consumer,
		SyncProducer:  syncProducer,
		AsyncProducer: aSyncProducer,
	}
	k.kafkaClient[topic] = newKafka
	return err
}

func (k kafka) CreateTopic(topic string, partitionNum int32) error {
	var err error = nil
	// Set number of partitions for topic
	err = k.admin.CreateTopic(
		topic, &sarama.TopicDetail{
			NumPartitions:     partitionNum,
			ReplicationFactor: 1,
		}, false,
	)

	if err != nil {
		//if strings.Contains(err.Error(), "already exists") {
		//	_ = k.admin.DeleteTopic(topic)
		//	return k.CreateTopic(topic, partitionNum)
		//}
		fmt.Println("Error setting number of partitions:=> ", err)
		os.Exit(1)
	}
	return err
}

func (k kafka) SyncProducerData(topic string, partition int32, content string) error {
	var err error = nil
	currentClient, isExists := k.kafkaClient[topic]
	if !isExists {
		return fmt.Errorf("kafka topic not found")
	}

	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: partition,
		Value:     sarama.StringEncoder(content),
	}
	mu := sync.Mutex{}

	mu.Lock()
	_, _, err = currentClient.SyncProducer.SendMessage(msg)
	mu.Unlock()
	if err != nil {
		fmt.Println("Failed to send message Sync", err)
		return err
	}

	return err
}

func (k kafka) ASyncProducerData(topic string, partition int32, content string) error {
	var err error = nil
	currentClient, isExists := k.kafkaClient[topic]
	if !isExists {
		return fmt.Errorf("kafka topic not found")
	}

	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: partition,
		Value:     sarama.StringEncoder(content),
	}

loop:
	for {
		currentClient.AsyncProducer.Input() <- msg

		select {
		case _ = <-currentClient.AsyncProducer.Successes():
			break loop
		case errorInfo := <-currentClient.AsyncProducer.Errors():
			err = errorInfo
			fmt.Println("Failed to send message Async:", err)
			break loop
		}
	}
	return err
}

func (k kafka) ConsumerData(topic string, partition int32, parse ParseStruct) ([]string, error) {
	var err error = nil
	currentClient, isExists := k.kafkaClient[topic]
	if !isExists {
		return nil, fmt.Errorf("kafka topic not found")
	}

	resp, err := currentClient.Consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if err != nil {
		return nil, err
	}
	//size := resp.HighWaterMarkOffset()
	//collection := make(chan sarama.ConsumerMessage, 100)
	var collection []sarama.ConsumerMessage
	//wg := sync.WaitGroup{}
	//wg.Add(1)
	//count := 0
loopConsumer:
	for {
		//fmt.Println("Consumer Info =>", topic, size, currentClient.Consumer.HighWaterMarks())
		select {
		case errInfo := <-resp.Errors():
			err = errInfo
			fmt.Println(err)
			break loopConsumer
		case msg := <-resp.Messages():
			//fmt.Println(string(msg.Value))
			if string(msg.Value) == model.EOF {
				break loopConsumer
			}
			collection = append(collection, *msg)
			//collection <- *msg
		}
	}
	//close(collection)

	if err != nil {
		return nil, err
	}
	result, err := parse(&collection)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (k kafka) CloseTopic() error {
	for _, client := range k.kafkaClient {
		if err := client.Consumer.Close(); err != nil {
			return err
		}
		if err := client.SyncProducer.Close(); err != nil {
			return err
		}
		if err := client.AsyncProducer.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (k kafka) RemoveTopic() error {
	defer func() { _ = k.admin.Close() }()
	for topic := range k.system.Topics {
		_ = k.admin.DeleteTopic(topic)
		//if err != nil {
		//	return nil
		//}
	}
	return nil
}
