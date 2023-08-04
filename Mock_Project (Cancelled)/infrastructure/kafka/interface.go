package kafka

import (
	"Mock_Project/model"
	"github.com/IBM/sarama"
)

// ParseStruct parse record from kafka to struct
type ParseStruct func(msg *chan sarama.ConsumerMessage) ([]interface{}, error)

type IKafkaHandler interface {
	CreateTopic(topic string, partitionNum int32) error
	ProducerData(broker []string, topic string, partition int32, content string) error
	ConsumerData(broker []string, topic string, partition int32, parse ParseStruct) ([]interface{}, error)
	ClearTopic(server *model.Server) error
}
