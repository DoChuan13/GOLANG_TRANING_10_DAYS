package kafka_process

import (
	"Mock_Project/infrastructure/repository"
	"Mock_Project/model"
	"Mock_Project/pkg"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type Server struct {
	config          *model.Server
	kafkaRepository repository.IKafkaRepository
	wg              *sync.WaitGroup
	mu              *sync.Mutex
	cond            *sync.Cond
	err             chan error
	producerLtd     chan bool
	consumerLtd     chan bool
	startTime       time.Time
}

var totalCount, reCount = 0, 0

func NewKafkaService(startTime time.Time, cfg *model.Server, kafkaRepository *repository.IKafkaRepository) IKafka {
	return &Server{
		config:          cfg,
		kafkaRepository: *kafkaRepository,
		wg:              new(sync.WaitGroup),
		mu:              new(sync.Mutex),
		cond:            new(sync.Cond),
		err:             make(chan error, 1),
		producerLtd:     make(chan bool, cfg.ProducerLtd),
		consumerLtd:     make(chan bool, cfg.ConsumerLtd),
		startTime:       startTime,
	}
}

func (s Server) StartKafkaProcess(csCh chan model.ConsumerObject, done chan bool, rows []string) error {
	var err error = nil
	pkg.LogStepProcess(s.startTime, "2.1 Start Producer")
	s.cond.L = new(sync.Mutex)

	//Clear All Topic Error
	//s.clearTopicAndStop(rows)

	//Init Kafka Connection
	err = s.kafkaRepository.InitConnection()
	if err != nil {
		return err
	}

	var recordSlide []string
	//Producer All Messages (Topic = First Character + Last Character)
	for i := 0; i < len(rows); i++ {
		s.producerLtd <- true
		totalCount++
		recordSlide = strings.Split(rows[i], model.CommaCharacter)
		topicName := recordSlide[2] + model.UnderScoreCharacter + recordSlide[3] + model.UnderScoreCharacter + recordSlide[4]
		_, isExists := s.config.Topics[topicName]

		//Init Topic If not Exist
		if !isExists {
			err = s.kafkaRepository.CreateTopic(topicName, s.config.MaxPartition)
			if err != nil {
				return err
			}
		}
		s.config.Topics[topicName]++

		//Limit Lock Number Goroutine
		s.cond.L.Lock()
		if len(s.producerLtd) == s.config.ProducerLtd {
			s.cond.Wait()
		}
		s.wg.Add(1)
		go s.producerProcess(topicName, rows[i], 0)
		s.cond.L.Unlock()
	}

	s.wg.Wait()
	close(s.producerLtd)

	//Detect Error in Goroutine
	err = s.breakError()
	if err != nil {
		fmt.Println("Error Producer==> ", err)
		return err
	}
	pkg.LogStepProcess(s.startTime, "2.2 Finish Producer")
	fmt.Println("Total GoRoutine of Producer ====> ", totalCount)

	pkg.LogStepProcess(s.startTime, "2.3 Start Consumer")
	for topic := range s.config.Topics {
		s.consumerLtd <- true
		s.cond.L.Lock()
		if len(s.consumerLtd) == s.config.ConsumerLtd {
			s.cond.Wait()
		}
		s.wg.Add(1)
		go s.consumerProcess(csCh, topic, 0)
		s.cond.L.Unlock()
	}

	s.wg.Wait()
	close(s.consumerLtd)

	//Detect Error in Goroutine
	err = s.breakError()
	if err != nil {
		return err
	}

	pkg.LogStepProcess(s.startTime, "2.4 Finish Consumer")

	done <- true
	return nil
}

func (s Server) breakError() error {
	select {
	case err := <-s.err:
		return err
	default:
	}
	return nil
}

func (s Server) producerProcess(topic, content string, partitionId int32) {
	err := s.kafkaRepository.SyncProducerData(topic, partitionId, content)
	//fmt.Println("Topic ===>", topic)
	if err != nil {
		s.err <- err
		return
	}
	s.wg.Done()
	<-s.producerLtd
	s.cond.Signal()
}

func (s Server) consumerProcess(csCha chan model.ConsumerObject, topic string, partitionId int32) {
	defer s.wg.Done()
	messages, err := s.kafkaRepository.ConsumerData(topic, partitionId)
	if err != nil {
		s.err <- err
		return
	}
	<-s.consumerLtd
	s.cond.Signal()
	reCount += len(messages)
	consumer := model.ConsumerObject{TableName: topic, Records: messages}

	csCha <- consumer
}

func (s Server) clearTopicAndStop(rows []string) {
	var recordSlide []string
	for i := 0; i < len(rows); i++ {
		recordSlide = strings.Split(rows[i], model.CommaCharacter)
		topicName := recordSlide[2] + model.UnderScoreCharacter + recordSlide[3] + model.UnderScoreCharacter + recordSlide[4]
		s.config.Topics[topicName]++
	}

	err := s.kafkaRepository.RemoveTopic()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Topic Clear Success")
	os.Exit(1)
}
