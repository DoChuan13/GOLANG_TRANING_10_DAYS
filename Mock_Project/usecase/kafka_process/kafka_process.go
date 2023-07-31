package kafka_process

import (
	"Mock_Project/infrastructure/repository"
	"Mock_Project/model"
	"encoding/json"
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
}

func NewService(cfg *model.Server, kafkaRepository *repository.IKafkaRepository) IKafka {
	return &Server{
		config:          cfg,
		kafkaRepository: *kafkaRepository,
		wg:              new(sync.WaitGroup),
		mu:              new(sync.Mutex),
	}
}

func (s Server) StartKafkaProcess(rows []string) ([]model.ObjectProcess, error) {
	//Producer All Messages (Topic = First Character + Last Character)
	for _, row := range rows {
		splitRow := strings.Split(row, model.CommaCharacter)
		topic := splitRow[2] + splitRow[3]
		s.saveListTopic(topic)
		s.wg.Add(1)
		go func(topic, content string) {
			_ = s.kafkaRepository.ProducerData(s.config.Broker, topic, s.config.Partition, content)
			s.wg.Done()
		}(topic, row)
	}
	s.wg.Wait()

	//Consumer All Messages & Return All Data
	var collection []model.ObjectProcess
	for _, topic := range s.config.Topics {
		s.wg.Add(1)

		go func(topic string, collection *[]model.ObjectProcess) {
			values, err := s.kafkaRepository.ConsumerData(s.config.Broker, topic, s.config.Partition)
			if err != nil {
				return
			}
			s.mergeTableGroup(collection, values)
			s.wg.Done()
		}(topic, &collection)

	}
	s.wg.Wait()

	result := s.sortValueByTime(collection)
	for i := 0; i < len(result); i++ {
		for j := 0; j < len(result[i].Value); j++ {
			curDate := formatDate(time.Now())
			curTime := formatTime(time.Now())
			result[i].Value[j].TKSERIALNUMBER = 1
			result[i].Value[j].TKZXD = curDate
			result[i].Value[j].TKTIM = curTime
			if j > 0 {
				if curDate == result[i].Value[j-1].TKZXD && curTime == result[i].Value[j-1].TKTIM {
					result[i].Value[j].TKSERIALNUMBER = result[i].Value[j-1].TKSERIALNUMBER + 1
				}
			}
		}
	}

	err := s.kafkaRepository.ClearData(s.config)
	if err != nil {
		return nil, err
	}

	return result, nil
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
				EndPoint:  "3306",
				DBName:    "",
				Value:     []model.TargetObject{targetObject},
			}
			objectList = append(objectList, object)
		}
	}
	*collection = objectList
	s.mu.Unlock()
}

func (s Server) sortValueByTime(collection []model.ObjectProcess) []model.ObjectProcess {
	for i := 0; i < len(collection); i++ {
		s.wg.Add(1)
		go s.sortItem(&collection[i])
	}

	s.wg.Wait()
	return collection
}

func (s Server) sortItem(collect *model.ObjectProcess) {
	compare := func(i, j int) bool {
		return compare(i, j, collect.Value)
	}
	sort.Slice(collect.Value, compare)
	s.wg.Done()
}

func compare(i, j int, value []model.TargetObject) bool {
	if value[i].QCD != value[j].QCD {
		return strings.Compare(value[i].QCD, value[j].QCD) < 0
	}
	return strings.Compare(value[i].TIME, value[j].TIME) < 0
}

func convertToObject(value []string) (tableName string, targetObject model.TargetObject) {
	table := value[2] + model.UnderScoreCharacter + value[3] + model.UnderScoreCharacter + value[4]
	tempObject := model.SourceObject{}
	temp, err := json.Marshal(tempObject)
	if err != nil {
		return
	}
	tempStr := strings.Trim(string(temp), "{}")
	nameArray := strings.Split(tempStr, ":\"\",")
	objectMap := make(map[string]string)

	for index, name := range nameArray {
		key := strings.Trim(name, "\"")
		objectMap[key] = value[index]
	}

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

func formatTime(time time.Time) string {
	return time.Format(model.TimeFormatWithMicrosecond)
}

func formatDate(time time.Time) string {
	return time.Format(model.DateFormatWithStroke)
}
