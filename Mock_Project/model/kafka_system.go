package model

type KafkaSystem struct {
	Broker    []string
	Topics    []string
	Partition int32
}
