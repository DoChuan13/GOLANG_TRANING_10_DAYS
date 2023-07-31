package kafka_process

import "Mock_Project/model"

type IKafka interface {
	StartKafkaProcess(rows []string) ([]model.ObjectProcess, error)
}
