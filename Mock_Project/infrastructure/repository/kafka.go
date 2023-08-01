package repository

import (
	"Mock_Project/infrastructure"
	"Mock_Project/infrastructure/kafka"
	"Mock_Project/model"
	"github.com/IBM/sarama"
)

type kafkaRepository struct {
	config      *model.Server
	kafkaClient kafka.IKafkaHandler
}

// NewKafkaRepository repository constructor
func NewKafkaRepository(infra *infrastructure.Infra, cfg *model.Server) IKafkaRepository {
	return &kafkaRepository{
		config:      cfg,
		kafkaClient: infra.KafkaHandler,
	}
}

func (k kafkaRepository) ProducerData(broker []string, topic string, partition int32, content string) error {
	return k.kafkaClient.ProducerData(broker, topic, partition, content)
}

func (k kafkaRepository) ConsumerData(broker []string, topic string, partition int32) ([]interface{}, error) {
	value, err := k.kafkaClient.ConsumerData(broker, topic, partition, parseConsumer)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (k kafkaRepository) ClearData(server *model.Server) error {
	return k.kafkaClient.ClearTopic(server)
}

// Parse ConsumerMessage to Slice
func parseConsumer(collection *chan sarama.ConsumerMessage) ([]interface{}, error) {
	var result []interface{}
	for msg := range *collection {
		result = append(result, string(msg.Value))
	}
	return result, nil
}
