package test

import (
	kafka2 "Mock_Project/infrastructure/kafka"
	"Mock_Project/model"
	"testing"
)

func TestKafka(t *testing.T) {
	cfg := model.KafkaSystem{
		Broker:    []string{"0.0.0.0:9093"},
		Topics:    "SZN-TSE1",
		Partition: 0,
	}
	kafka, err := kafka2.NewKafkaHandler(&cfg)
	if err != nil {
		return
	}
	_ = kafka.ConsumerData(cfg.Broker, cfg.Topics, cfg.Partition)
}
