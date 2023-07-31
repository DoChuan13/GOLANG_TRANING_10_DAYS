package main

import (
	"Mock_Project/infrastructure/repository"
	"Mock_Project/model"
	"runtime"
)

func init() {
	initLog()
}

func initLog() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

var pathFile = "file/demo.csv"

func main() {
	temp := model.TargetObject{
		TKTIM: "Hello",
	}
	temphe := model.TargetObject{
		TKTIM: "World",
	}
	temp2 := []model.TargetObject{temp, temphe}
	_ = repository.ConvertColumns(temp)
	_ = repository.ConvertValues(temp2)
	//fmt.Println(colums, values)
	//fileService := read_data.NewService()
	//row, err := fileService.ReadFileProcess(pathFile)
	//if err != nil {
	//	return
	//}
	////fmt.Println(row)
	//
	//kafkaDB := model.KafkaDB{}
	//kafkaSystem := model.KafkaSystem{
	//	Broker:    []string{"0.0.0.0:9093"},
	//	Topics:    []string{},
	//	Partition: 0,
	//}
	//cfg := model.Server{
	//	KafkaDB:     kafkaDB,
	//	KafkaSystem: kafkaSystem,
	//}
	//kafkaHandler, _ := kafka.NewKafkaHandler(&kafkaSystem)
	//kafkaRepository := repository.NewKafkaRepository(kafkaHandler)
	//kafkaProcess := kafka_process.NewService(&cfg, &kafkaRepository)
	//
	//_, _ = kafkaProcess.StartKafkaProcess(row)
}
