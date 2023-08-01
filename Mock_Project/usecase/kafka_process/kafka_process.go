package kafka_process

import (
	"Mock_Project/infrastructure/repository"
	"Mock_Project/model"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"
)

type Server struct {
	config          *model.Server
	kafkaRepository repository.IKafkaRepository
	wg              *sync.WaitGroup
	mu              *sync.Mutex
	err             chan error
}

var totalCount = 0

func NewKafkaService(cfg *model.Server, kafkaRepository *repository.IKafkaRepository) IKafka {
	return &Server{
		config:          cfg,
		kafkaRepository: *kafkaRepository,
		wg:              new(sync.WaitGroup),
		mu:              new(sync.Mutex),
		err:             make(chan error, 1),
	}
}

func (s Server) StartKafkaProcess(rows []string) ([]model.ObjectProcess, error) {
	countElement := 0
	//Producer All Messages (Topic = First Character + Last Character)
	groupData := processGroupTopic(rows)
	for _, kafkaObj := range groupData {
		block := 1000
		size := (len(kafkaObj.ListRows) + block - 1) / block
		countElement += size

		for i := 0; i < size; i++ {
			start := i * block
			end := (i + 1) * block
			if end > len(kafkaObj.ListRows) {
				end = len(kafkaObj.ListRows)
			}
			batch := kafkaObj.ListRows[start:end]
			message := buildMessage(batch)
			topic := kafkaObj.Topic
			s.saveListTopic(topic)
			s.wg.Add(1)
			go s.producerProcess(topic, message)
		}

	}

	s.wg.Wait()

	defer func() {
		_ = s.kafkaRepository.ClearData(s.config)
	}()

	isProduced := false
	for !isProduced {
		select {
		case err := <-s.err:
			return nil, err
		default:
			isProduced = true
		}
	}

	//Consumer All Messages & Return All Data
	var collection []model.ObjectProcess
	for _, topic := range s.config.Topics {
		s.wg.Add(1)
		go s.consumerProcess(topic, &collection)
	}
	s.wg.Wait()

	count := 0
	for _, collect := range collection {
		count += len(collect.Value)
	}
	fmt.Println("Consumer Data ====> ", count)

	result := s.prepareDBList(collection)
	return result, nil
}

func buildMessage(batch []string) string {
	return strings.Join(batch, model.NewLineCharacter)
}

func processGroupTopic(rows []string) []model.KafkaProcess {
	var kafkaObjectList []model.KafkaProcess
	for _, row := range rows {
		temp := strings.Split(fmt.Sprint(row), model.CommaCharacter)
		topic := temp[2] + temp[3]
		existed := false
		index := -1
		for i := 0; i < len(kafkaObjectList); i++ {
			if kafkaObjectList[i].Topic == topic {
				existed = true
				index = i
				break
			}
		}
		if existed {
			kafkaObjectList[index].ListRows = append(kafkaObjectList[index].ListRows, row)
		} else {
			newKafkaObject := model.KafkaProcess{
				Topic:    topic,
				Message:  "",
				ListRows: []string{row},
			}
			kafkaObjectList = append(kafkaObjectList, newKafkaObject)
		}
	}
	return kafkaObjectList
}

func (s Server) saveListTopic(newTopic string) {
	existed := false
	for i := 0; i < len(s.config.Topics); i++ {
		if s.config.Topics[i] == newTopic {
			existed = true
			break
		}
	}
	if !existed {
		s.config.Topics = append(s.config.Topics, newTopic)
	}
}

func (s Server) producerProcess(topic string, content string) {
	err := s.kafkaRepository.ProducerData(s.config.Broker, topic, s.config.Partition, content)
	if err != nil {
		s.err <- err
		fmt.Println("Producer Error =>: ", err)
	}
	s.mu.Lock()
	totalCount++
	s.mu.Unlock()
	s.wg.Done()
}

func (s Server) consumerProcess(topic string, collection *[]model.ObjectProcess) {
	defer s.wg.Done()
	messages, err := s.kafkaRepository.ConsumerData(s.config.Broker, topic, s.config.Partition)
	if err != nil {
		fmt.Println("Consumer Error =>: ", err)
		return
	}
	s.mergeTableGroup(collection, messages)
}

func (s Server) mergeTableGroup(collection *[]model.ObjectProcess, messages []interface{}) {
	s.mu.Lock()
	objectList := *collection
	for _, message := range messages {
		var row = strings.Split(
			strings.ReplaceAll(fmt.Sprint(message), model.DoubleQuote, model.EmptyString), model.NewLineCharacter,
		)
		for _, value := range row {
			tableName, _, targetObject := convertToObject(fmt.Sprint(value))
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
					EndPoint:  "127.0.0.1",
					DBName:    "demo",
					Value:     []model.TargetObject{targetObject},
				}
				objectList = append(objectList, object)
			}
		}
	}

	*collection = objectList
	s.mu.Unlock()
}

func (s Server) prepareDBList(collection []model.ObjectProcess) []model.ObjectProcess {
	for i := 0; i < len(collection); i++ {
		s.wg.Add(1)
		go s.sortItems(&collection[i])
	}
	s.wg.Wait()
	for i := 0; i < len(collection); i++ {
		for j := 0; j < len(collection[i].Value); j++ {
			curDate := formatDate(time.Now())
			curTime := formatTime(time.Now())
			collection[i].Value[j].TKSERIALNUMBER = 1
			collection[i].Value[j].TKZXD = curDate
			collection[i].Value[j].TKTIM = curTime
			if j > 0 {
				if collection[i].Value[j].QCD == collection[i].Value[j-1].QCD && curDate == collection[i].Value[j-1].TKZXD && curTime == collection[i].Value[j-1].TKTIM {
					collection[i].Value[j].TKSERIALNUMBER = collection[i].Value[j-1].TKSERIALNUMBER + 1
				}
			}
		}
	}
	return collection
}

func (s Server) sortItems(collect *model.ObjectProcess) {
	defer s.wg.Done()
	compare := func(i, j int) bool {
		return compareObject(i, j, collect.Value)
	}
	sort.Slice(collect.Value, compare)
}

func compareObject(i, j int, value []model.TargetObject) bool {
	if value[i].QCD != value[j].QCD {
		return strings.Compare(value[i].QCD, value[j].QCD) < 0
	}
	return strings.Compare(value[i].TIME, value[j].TIME) < 0
}

func convertToObject(valueStr string) (tableName, topicName string, targetObject model.TargetObject) {
	value := strings.Split(fmt.Sprint(valueStr), model.CommaCharacter)
	table := value[2] + model.UnderScoreCharacter + value[3] + model.UnderScoreCharacter + value[4]
	topic := value[2] + value[3]
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
	return table, topic, object
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

func formatTime(time time.Time) string {
	return time.Format(model.TimeFormatWithMicrosecond)
}

func formatDate(time time.Time) string {
	return time.Format(model.DateFormatWithStroke)
}
