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
	ltdRoutine      chan bool
}

var totalCount, reCount = 0, 0

func NewKafkaService(cfg *model.Server, kafkaRepository *repository.IKafkaRepository) IKafka {
	return &Server{
		config:          cfg,
		kafkaRepository: *kafkaRepository,
		wg:              new(sync.WaitGroup),
		mu:              new(sync.Mutex),
		cond:            new(sync.Cond),
		err:             make(chan error, 1),
		ltdRoutine:      make(chan bool, cfg.Limited),
	}
}

var startTime = time.Now()

func (s Server) StartKafkaProcess(csCh chan model.ConsumerObject, done chan bool, rows []string) error {
	pkg.LogStepProcess(startTime, "2.1 Start Producer")
	s.cond.L = new(sync.Mutex)

	//Clear All Topic Error
	//s.clearTopicAndStop(rows)

	var recordSlide []string
	//Producer All Messages (Topic = First Character + Last Character)
	for i := 0; i < len(rows); i++ {
		s.ltdRoutine <- true
		totalCount++
		recordSlide = strings.Split(rows[i], model.CommaCharacter)
		topicName := recordSlide[2] + model.UnderScoreCharacter + recordSlide[3] + model.UnderScoreCharacter + recordSlide[4]
		_, isExists := s.config.Topics[topicName]

		//Init Kafka Connect If not Exist
		if !isExists {
			err := s.kafkaRepository.InitConnection(topicName)
			if err != nil {
				return err
			}
		}
		s.config.Topics[topicName]++

		//Limit Number Goroutine
		s.cond.L.Lock()
		if len(s.ltdRoutine) == s.config.Limited {
			//fmt.Println("Waiting Limited Goroutine")
			s.cond.Wait()
		}
		s.wg.Add(1)
		go s.producerProcess(topicName, rows[i], 0)
		s.cond.L.Unlock()
	}

	s.wg.Wait()
	close(s.ltdRoutine)

	//Detect Error in Goroutine
	err := s.breakError()
	if err != nil {
		fmt.Println("Error Producer==> ", err)
		return err
	}
	pkg.LogStepProcess(startTime, "2.2 Finish Producer")
	fmt.Println("Total GoRoutine of Producer ====> ", totalCount)

	pkg.LogStepProcess(startTime, "2.3 Start Consumer")
	for topic := range s.config.Topics {
		s.wg.Add(1)
		go s.consumerProcess(csCh, topic, 0)
	}

	s.wg.Wait()

	//Detect Error in Goroutine
	err = s.breakError()
	if err != nil {
		return err
	}

	pkg.LogStepProcess(startTime, "2.4 Finish Consumer")
	//fmt.Println("Re-check Collection (Total Consumer Data)... ", reCount)

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
	<-s.ltdRoutine
	s.cond.Signal()
}

func (s Server) consumerProcess(csCha chan model.ConsumerObject, topic string, partitionId int32) {
	defer s.wg.Done()
	messages, err := s.kafkaRepository.ConsumerData(topic, partitionId)
	if err != nil {
		s.err <- err
		return
	}
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
