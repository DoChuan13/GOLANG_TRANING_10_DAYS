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
	clientName  string
	err         chan error
	wg          *sync.WaitGroup
}

// NewKafkaHandler constructor
func NewKafkaHandler(cfg *model.Server) (IKafkaHandler, error) {
	config := sarama.NewConfig()

	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.MaxMessageBytes = 10e6
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	config.Consumer.Offsets.Initial = sarama.OffsetOldest
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
		clientName:  "KafkaClient",
		err:         make(chan error, 1),
		wg:          &sync.WaitGroup{},
	}, err
}

func (k kafka) InitConnection() error {
	var err error = nil
	_, isExists := k.kafkaClient[k.clientName]
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
	k.kafkaClient[k.clientName] = newKafka
	return err
}

func (k kafka) CreateTopic(topic string, partitionNum int32) error {
	var err error = nil
	//topicList, err := k.admin.ListTopics()
	//if err != nil {
	//	return err
	//}
	//_, exist := topicList[topic]
	//if exist {
	//	err = k.admin.DeleteTopic(topic)
	//	if err != nil {
	//		return err
	//	}
	//	return k.CreateTopic(topic, partitionNum)
	//}
	// Set number of partitions for topic
	err = k.admin.CreateTopic(
		topic, &sarama.TopicDetail{
			NumPartitions:     partitionNum,
			ReplicationFactor: 1,
		}, false,
	)

	if err != nil {
		fmt.Println("Error setting number of partitions:=> ", err)
		os.Exit(1)
	}
	return err
}

func (k kafka) SyncProducerData(topic string, partition int32, content string) error {
	var err error = nil
	currentClient, isExists := k.kafkaClient[k.clientName]
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
	currentClient, isExists := k.kafkaClient[k.clientName]
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

func (k kafka) ConsumerData(topic string, _ int32, parse ParseStruct) ([]string, error) {
	var err error = nil
	currentClient, isExists := k.kafkaClient[k.clientName]
	if !isExists {
		return nil, fmt.Errorf("kafka topic not found")
	}

	collection := make(chan sarama.ConsumerMessage, k.system.Topics[topic])
	for i := 0; i < int(k.system.MaxPartition); i++ {
		k.wg.Add(1)
		go k.processConsumerPartitions(topic, currentClient, int32(i), collection)
	}
	k.wg.Wait()
	close(collection)

	//Detect Error in Goroutine
	err = k.breakError()
	if err != nil {
		fmt.Println("Error Producer==> ", err)
		return nil, err
	}
	result, err := parse(&collection)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (k kafka) processConsumerPartitions(
	topic string, kafka model.Kafka, partitionId int32, collection chan sarama.ConsumerMessage,
) {
	count := 0
	resp, err := kafka.Consumer.ConsumePartition(topic, partitionId, sarama.OffsetOldest)
	if err != nil {
		k.err <- err
		return
	}
	size := resp.HighWaterMarkOffset()
	isContinue := true
	for isContinue {
		select {
		case errInfo := <-resp.Errors():
			err = errInfo
			fmt.Println("Consumer Partition err ==> ", err)
			isContinue = false
		case msg := <-resp.Messages():
			collection <- *msg
			count++
			if count == int(size) {
				isContinue = false
			}
		}
	}
	k.wg.Done()
}

func (k kafka) breakError() error {
	select {
	case err := <-k.err:
		return err
	default:
	}
	return nil
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
	topicList, err := k.admin.ListTopics()
	if err != nil {
		return err
	}
	for topic := range topicList {
		_ = k.admin.DeleteTopic(topic)
	}
	return nil
}
