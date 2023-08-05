package kafka

import (
	"github.com/IBM/sarama"
)

// ParseStruct parse record from kafka to struct
type ParseStruct func(msg *chan sarama.ConsumerMessage) ([]string, error)

type IKafkaHandler interface {
	InitConnection(topic string) error
	CreateTopic(topic string, partitionNum int32) error
	SyncProducerData(topic string, partition int32, content string) error
	ASyncProducerData(topic string, partition int32, content string) error
	ConsumerData(topic string, partition int32, parse ParseStruct) ([]string, error)
	CloseTopic() error
	RemoveTopic() error
}
