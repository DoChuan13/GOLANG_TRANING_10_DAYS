package model

import "github.com/IBM/sarama"

type KafkaSystem struct {
	Broker       []string
	Topics       map[string]int
	MaxPartition int32
	Block        int
}

type Topic struct {
	Name string
}

type Kafka struct {
	Consumer sarama.Consumer
	Producer sarama.SyncProducer
}
