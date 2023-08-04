package kafka

import (
	"github.com/IBM/sarama"
)

// ParseStruct parse record from kafka to struct
type ParseStruct func(msg *chan sarama.ConsumerMessage) ([]string, error)

type IKafkaHandler interface {
	CreateTopic(topic string, partitionNum int32) error
	ProducerData(topic string, partition int32, content string) error
	ConsumerData(topic string, partition int32, parse ParseStruct) ([]string, error)
	InitConnection(topic string) error
	CloseTopic() error
	RemoveTopic() error
}
