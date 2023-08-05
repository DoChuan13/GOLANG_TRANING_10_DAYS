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

func (k kafkaRepository) SyncProducerData(topic string, partition int32, content string) error {
	return k.kafkaClient.SyncProducerData(topic, partition, content)
}

func (k kafkaRepository) ASyncProducerData(topic string, partition int32, content string) error {
	return k.kafkaClient.ASyncProducerData(topic, partition, content)
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
func parseConsumer(collection *[]sarama.ConsumerMessage) ([]string, error) {
	var result []string
	i := 0
	var newField [3]string
	for _, msg := range *collection {
		recordSlice := strings.Split(string(msg.Value), model.CommaCharacter)
		curDate := formatDate(time.Now())
		curTime := formatTime(time.Now())
		serialNumber := 1

		if i > 0 {
			if curDate == newField[0] && curTime == newField[1] {
				curSerial, _ := strconv.Atoi(newField[2])
				serialNumber = curSerial + 1
			}
		}
		newField[0] = curDate
		newField[1] = curTime
		newField[2] = strconv.Itoa(serialNumber)
		i++

		newField := []string{curDate, curTime, strconv.Itoa(serialNumber)}
		recordSlice = append(recordSlice, newField...)
		copy(recordSlice[2+len(newField):], recordSlice[2:])
		copy(recordSlice[2:], newField)
		result = append(result, strings.Join(recordSlice, model.CommaCharacter))
	}
	//close(*collection)
	return result, nil
}

func formatTime(time time.Time) string {
	return time.Format(model.TimeFormatWithMicrosecond)
}

func formatDate(time time.Time) string {
	return time.Format(model.DateFormatWithStroke)
}
