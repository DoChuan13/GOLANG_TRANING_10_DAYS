package model

type KafkaSystem struct {
	Broker       []string
	Topics       []Topic
	MaxPartition int32
	Block        int
}

type Topic struct {
	Name    string
	ParSize int32
}
