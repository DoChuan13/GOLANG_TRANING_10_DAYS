package api

import (
	"Mock_Project/infrastructure"
	"Mock_Project/infrastructure/repository"
	"Mock_Project/model"
	"Mock_Project/pkg/logger"
	"Mock_Project/usecase/fetch_db"
	"Mock_Project/usecase/insert_data"
	"Mock_Project/usecase/kafka_process"
	"Mock_Project/usecase/read_data"
	"context"
	"fmt"
	"time"
)

type Server struct {
	infra *infrastructure.Infra
	cfg   *model.Server
}

// New create new server instance
func New(infra *infrastructure.Infra, cfg *model.Server) *Server {
	return &Server{
		infra: infra,
		cfg:   cfg,
	}
}

func (s *Server) Start(ctx context.Context) error {
	//Config
	var mkdirPath = "file"
	var fileName = "ListValue.csv"
	folder := "/Users/Chuan/Desktop/demo/temp/"
	kafkaDB := model.KafkaDB{
		Port:              3306,
		User:              "root",
		Password:          "ChuanDo@13",
		MaxOpenConnection: 10,
		MaxIdleConnection: 10,
		DriverName:        "mysql",
		RetryTimes:        3000,
		RetryWaitMs:       3000,
	}
	kafkaSystem := model.KafkaSystem{
		Broker:    []string{"0.0.0.0:9093"},
		Topics:    []string{},
		Partition: 0,
	}
	cfg := model.Server{
		KafkaDB:     kafkaDB,
		KafkaSystem: kafkaSystem,
	}

	//Read File Process
	getLogInfo("Step 1: Read File Started \t\t(%s)\n")
	fileService := read_data.NewService(mkdirPath, fileName)
	rows, err := fileService.ReadFileProcess()
	if err != nil {
		return err
	}
	getLogInfo("Step 1: Read File Completed \t\t(%s)\n")

	//Load Kafka Repository
	kafkaRepository := repository.NewKafkaRepository(s.infra, &cfg)

	//Load Database Repository
	dbRepository := repository.NewDBRepository(s.infra, &cfg)

	//Start Kafka Service & Process
	getLogInfo("Step 2: Kafka Process Started \t\t(%s)\n")
	kafkaService := kafka_process.NewKafkaService(&cfg, &kafkaRepository)
	objectProcessList, err := kafkaService.StartKafkaProcess(rows)
	if err != nil {
		return err
	}
	getLogInfo("Step 2: Kafka Process Completed \t(%s)\n")

	//Start Fetch Database Service & Process
	getLogInfo("Step 3: Fetch DB Started \t\t(%s)\n")
	fetchService := fetch_db.NewFetchService(&cfg, &dbRepository)
	err = fetchService.StartFetchDB(ctx, &objectProcessList)
	if err != nil {
		return err
	}
	getLogInfo("Step 3: Fetch DB Completed \t\t(%s)\n")

	//Start Database Service & Process
	getLogInfo("Step 4: Import DB Started \t\t(%s)\n")
	dbService := insert_data.NewDBService(&cfg, &dbRepository, folder)
	err = dbService.StartDBProcess(ctx, &objectProcessList)
	if err != nil {
		return err
	}
	getLogInfo("Step 4: Import DB Completed \t\t(%s)\n")

	getLogInfo("Process Finished at \t\t\t(%s)\n")
	return nil
}

// Close closes server and related resources
func (s *Server) Close() {
	s.infra.Close()
	logger.AppLog(logger.InfoAppNewRecoverDataFinished, logger.InfoLevelLog)
}

func getLogInfo(log string) {
	fmt.Printf(log, time.Now().Format(time.StampMicro))
}
