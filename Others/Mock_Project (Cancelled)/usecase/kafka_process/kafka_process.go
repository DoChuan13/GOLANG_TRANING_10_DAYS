package kafka_process

import (
	"Mock_Project/infrastructure/repository"
	"Mock_Project/model"
	"Mock_Project/pkg"
	"fmt"
	"reflect"
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
	message         chan model.ConsumerObject
	ltdRoutine      chan bool
}

var totalCount = 0

func NewKafkaService(cfg *model.Server, kafkaRepository *repository.IKafkaRepository) IKafka {
	return &Server{
		config:          cfg,
		kafkaRepository: *kafkaRepository,
		wg:              new(sync.WaitGroup),
		mu:              new(sync.Mutex),
		cond:            new(sync.Cond),
		err:             make(chan error, 1),
		message:         make(chan model.ConsumerObject, 1),
		ltdRoutine:      make(chan bool, cfg.Limited),
	}
}

var startTime = time.Now()

func (s Server) StartKafkaProcess(rows []string) ([]model.ObjectProcess, error) {
	pkg.LogStepProcess(startTime, "2.1 Start Producer")
	countElement := 0
	s.cond.L = new(sync.Mutex)
	//Producer All Messages (Topic = First Character + Last Character)
	groupData := processGroupTopic(rows)
	for _, kafkaObj := range groupData {
		size := (len(kafkaObj.ListRows) + s.config.Block - 1) / s.config.Block
		countElement += size
		for i, j := 0, 1; i < size && j < int(s.config.MaxPartition); i, j = i+1, j+1 {
			start := i * s.config.Block
			end := (i + 1) * s.config.Block
			if end > len(kafkaObj.ListRows) {
				end = len(kafkaObj.ListRows)
			}
			batch := kafkaObj.ListRows[start:end]
			message := buildMessage(batch)
			topic := kafkaObj.Topic

			s.saveListTopic(topic)

			totalCount++
			s.cond.L.Lock()
			s.ltdRoutine <- true
			if len(s.ltdRoutine) == s.config.Limited {
				//fmt.Println("Waiting Limited Goroutine")
				s.cond.Wait()
			}
			//fmt.Println("Next Goroutine")
			s.wg.Add(1)
			go s.producerProcess(topic, message, int32(j))
			if j == int(s.config.MaxPartition)-1 {
				j = -1
			}
			s.cond.L.Unlock()
		}
	}
	s.wg.Wait()
	close(s.ltdRoutine) //Detect Error in Goroutine
	err := s.breakError()
	if err != nil {
		return nil, err
	}
	pkg.LogStepProcess(startTime, "2.2 Finish Producer")
	fmt.Println("Total GoRoutine of Producer ====> ", totalCount)

	defer func() {
		_ = s.kafkaRepository.ClearData(s.config)
	}()

	//Consumer All Messages & Return All Data
	var collection []model.ObjectProcess
	pkg.LogStepProcess(startTime, "2.3 Start Consumer")
	topics := s.config.Topics
	for i := 0; i < len(topics); i++ {
		for j := 0; j < int(s.config.MaxPartition); j++ {
			s.wg.Add(1)
			go s.consumerProcess(topics[i].Name, &collection, int32(i))
		}
	}
	s.wg.Wait()

	//Detect Error in Goroutine
	err = s.breakError()
	if err != nil {
		return nil, err
	}

	pkg.LogStepProcess(startTime, "2.4 Finish Consumer")
	fmt.Println("Re-check Collection... ")
	count := 0
	for _, collect := range collection {
		count += len(collect.Value)
	}
	fmt.Println("Total Consumer Data ====> ", count)
	return collection, nil
}

func (s Server) breakError() error {
	select {
	case err := <-s.err:
		return err
	default:
	}
	return nil
}

func buildMessage(batch []string) string {
	return strings.Join(batch, model.NewLineCharacter)
}

func processGroupTopic(rows []string) []model.KafkaProcess {
	var kafkaObjs []model.KafkaProcess
	for _, row := range rows {
		temp := strings.Split(fmt.Sprint(row), model.CommaCharacter)
		topic := temp[2] + temp[3]
		existed := false
		index := -1
		for i := 0; i < len(kafkaObjs); i++ {
			if kafkaObjs[i].Topic == topic {
				existed = true
				index = i
				break
			}
		}
		if existed {
			kafkaObjs[index].ListRows = append(kafkaObjs[index].ListRows, row)
		} else {
			newKafkaObject := model.KafkaProcess{
				Topic:    topic,
				ListRows: []string{row},
			}
			kafkaObjs = append(kafkaObjs, newKafkaObject)
		}
	}
	return kafkaObjs
}

func (s Server) saveListTopic(newTopic string) {
	existed := false
	for i := 0; i < len(s.config.Topics); i++ {
		if s.config.Topics[i].Name == newTopic {
			existed = true
			break
		}
	}

	if !existed {
		s.config.Topics = append(s.config.Topics, model.Topic{Name: newTopic, ParSize: s.config.MaxPartition})
		err := s.kafkaRepository.CreateTopic(newTopic, s.config.MaxPartition)
		if err != nil {
			s.err <- err
		}
	}
}

func (s Server) producerProcess(topic, content string, partitionId int32) {
	err := s.kafkaRepository.ProducerData(s.config.Broker, topic, partitionId, content)
	if err != nil {
		s.err <- err
		return
	}
	s.wg.Done()
	s.mu.Lock()
	<-s.ltdRoutine
	s.mu.Unlock()
	s.cond.Signal()
}

func (s Server) consumerProcess(topic string, collection *[]model.ObjectProcess, partitionId int32) {
	defer s.wg.Done()
	messages, err := s.kafkaRepository.ConsumerData(s.config.Broker, topic, partitionId)
	if err != nil {
		s.err <- err
		return
	}
	s.mergeTableGroup(collection, messages)
}

func (s Server) mergeTableGroup(collection *[]model.ObjectProcess, messages []interface{}) {
	s.mu.Lock()
	objectList := *collection
	for _, message := range messages {
		var row = strings.Split(fmt.Sprint(message), model.NewLineCharacter)
		for _, value := range row {
			tableName, targetObject := convertToObject(fmt.Sprint(value))
			existedTable := false
			index := -1
			for i := 0; i < len(objectList); i++ {
				if objectList[i].TableName == tableName {
					existedTable = true
					index = i
					break
				}
			}
			if existedTable {
				objectList[index].Value = append(objectList[index].Value, targetObject)
			} else {
				object := model.ObjectProcess{
					TableName: tableName,
					EndPoint:  s.config.Endpoint,
					DBName:    s.config.DBName,
					Value:     []model.TargetObject{targetObject},
				}
				objectList = append(objectList, object)
			}
		}
	}

	*collection = objectList
	s.mu.Unlock()
}

func convertToObject(valueStr string) (tableName string, targetObject model.TargetObject) {
	value := strings.Split(fmt.Sprint(valueStr), model.CommaCharacter)
	table := value[2] + model.UnderScoreCharacter + value[3] + model.UnderScoreCharacter + value[4]
	objectMap := generateObjectMap(value)

	object := model.TargetObject{}
	val := reflect.ValueOf(&object).Elem()
	typ := reflect.TypeOf(&object).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		for key, value := range objectMap {
			if field.Name == key {
				val.Field(i).SetString(value)
			}
		}
	}
	return table, object
}

func generateObjectMap(value []string) map[string]string {
	objectMap := make(map[string]string)
	var object model.SourceObject
	val := reflect.ValueOf(object)
	typ := reflect.TypeOf(object)
	for i := 0; i < val.NumField(); i++ {
		key := typ.Field(i).Name
		objectMap[key] = value[i]
	}
	return objectMap
}
