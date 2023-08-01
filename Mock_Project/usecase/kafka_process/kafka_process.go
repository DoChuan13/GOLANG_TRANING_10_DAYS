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
	collectCh       chan model.ObjectProcess
}

func NewKafkaService(cfg *model.Server, kafkaRepository *repository.IKafkaRepository) IKafka {
	return &Server{
		config:          cfg,
		kafkaRepository: *kafkaRepository,
		wg:              new(sync.WaitGroup),
		mu:              new(sync.Mutex),
		collectCh:       make(chan model.ObjectProcess, 1),
	}
}

func (s Server) StartKafkaProcess(rows []string) ([]model.ObjectProcess, error) {
	//Producer All Messages (Topic = First Character + Last Character)
	for _, row := range rows {
		splitRow := strings.Split(row, model.CommaCharacter)
		topic := splitRow[2] + splitRow[3]
		s.saveListTopic(topic)
		s.wg.Add(1)
		go s.producerProcess(topic, row)
	}
	s.wg.Wait()

	//Consumer All Messages & Return All Data
	s.collectCh = make(chan model.ObjectProcess, len(s.config.Topics))
	var collection []model.ObjectProcess
	for _, topic := range s.config.Topics {
		s.wg.Add(1)
		go s.consumerProcess(topic, &collection)
	}
	s.wg.Wait()

	result := s.prepareDBList(collection)

	err := s.kafkaRepository.ClearData(s.config)
	if err != nil {
		return nil, err
	}

	return result, nil
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
	defer s.wg.Done()
	_ = s.kafkaRepository.ProducerData(s.config.Broker, topic, s.config.Partition, content)
}

func (s Server) consumerProcess(topic string, collection *[]model.ObjectProcess) {
	defer s.wg.Done()
	values, err := s.kafkaRepository.ConsumerData(s.config.Broker, topic, s.config.Partition)
	if err != nil {
		return
	}
	s.mergeTableGroup(collection, values)
}

func (s Server) mergeTableGroup(collection *[]model.ObjectProcess, row []interface{}) {
	s.mu.Lock()
	objectList := *collection
	for _, value := range row {
		splitValue := strings.Split(fmt.Sprint(value), model.CommaCharacter)
		tableName, targetObject := convertToObject(splitValue)
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

func convertToObject(value []string) (tableName string, targetObject model.TargetObject) {
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

func formatTime(time time.Time) string {
	return time.Format(model.TimeFormatWithMicrosecond)
}

func formatDate(time time.Time) string {
	return time.Format(model.DateFormatWithStroke)
}
