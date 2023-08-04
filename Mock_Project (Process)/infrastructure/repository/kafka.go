package repository

import (
	"Mock_Project/infrastructure"
	"Mock_Project/infrastructure/kafka"
	"Mock_Project/model"
	"github.com/IBM/sarama"
	"strconv"
	"strings"
	"time"
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

func (k kafkaRepository) InitConnection(topic string) error {
	return k.kafkaClient.InitConnection(topic)
}

func (k kafkaRepository) CloseTopic() error {
	return k.kafkaClient.CloseTopic()
}

func (k kafkaRepository) CreateTopic(topic string, partitionNum int32) error {
	return k.kafkaClient.CreateTopic(topic, partitionNum)
}

func (k kafkaRepository) ProducerData(topic string, partition int32, content string) error {
	return k.kafkaClient.ProducerData(topic, partition, content)
}

func (k kafkaRepository) ConsumerData(topic string, partition int32) ([]string, error) {
	value, err := k.kafkaClient.ConsumerData(topic, partition, parseConsumer)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (k kafkaRepository) RemoveTopic() error {
	return k.kafkaClient.RemoveTopic()
}

// Parse ConsumerMessage to Slice
func parseConsumer(collection *chan sarama.ConsumerMessage) ([]string, error) {
	var result []string
	i := 0
	var preDate, preTime string
	var preSerial int
	for msg := range *collection {
		tempSlice := strings.Split(string(msg.Value), model.CommaCharacter)
		curDate := formatDate(time.Now())
		curTime := formatTime(time.Now())
		serialNumber := 1

		if i > 0 {
			if curDate == preDate && curTime == preTime {
				serialNumber = preSerial + 1
			}
		}
		preDate = curDate
		preTime = curTime
		preSerial = serialNumber
		i++

		newField := []string{curDate, curTime, strconv.Itoa(serialNumber)}
		finalValue := append(tempSlice[:2], append(newField, tempSlice[2:]...)...)
		result = append(result, strings.Join(finalValue, model.CommaCharacter))
	}
	return result, nil
}

func formatTime(time time.Time) string {
	return time.Format(model.TimeFormatWithMicrosecond)
}

func formatDate(time time.Time) string {
	return time.Format(model.DateFormatWithStroke)
}
