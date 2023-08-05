package kafka_process

import "Mock_Project/model"

type IKafka interface {
	StartKafkaProcess(consumer chan model.ConsumerObject, done chan bool, rows []string) error
}
