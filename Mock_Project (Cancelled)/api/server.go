package api

import (
	"Mock_Project/infrastructure"
	"Mock_Project/infrastructure/repository"
	"Mock_Project/model"
	"Mock_Project/pkg"
	"Mock_Project/pkg/logger"
	"Mock_Project/usecase/fetch_db"
	"Mock_Project/usecase/insert_data"
	"Mock_Project/usecase/kafka_process"
	"Mock_Project/usecase/read_data"
	"context"
	"fmt"
	"os"
	_ "path/filepath"
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
	dir, err := os.Getwd()
	fmt.Println(dir)
	if err != nil {
		return err
	}
	var mkdirPath = "file"
	var fileName = "ListValue200.000.csv"
	tempTarget := dir + "/file/temp"

	//err = file.FakeAllData()
	//if err != nil {
	//	fmt.Println(err)
	//	return err
	//}

	var startTime = time.Now()
	//Read File Process
	pkg.LogStepProcess(startTime, "Step 1: Read File Started")
	fileService := read_data.NewService(mkdirPath, fileName)
	rows, err := fileService.ReadFileProcess()
	if err != nil {
		return err
	}
	pkg.LogStepProcess(startTime, "Step 1: Read File Completed")

	//Load Kafka Repository
	kafkaRepository := repository.NewKafkaRepository(s.infra, s.cfg)

	//Load Database Repository
	dbRepository := repository.NewDBRepository(s.infra, s.cfg)

	//Start Kafka Service & Process
	pkg.LogStepProcess(startTime, "Step 2: Kafka Process Started")
	kafkaService := kafka_process.NewKafkaService(s.cfg, &kafkaRepository)
	objectProcessList, err := kafkaService.StartKafkaProcess(rows)
	if err != nil {
		return err
	}
	pkg.LogStepProcess(startTime, "Step 2: Kafka Process Completed")

	//Start Fetch Database Service & Process
	pkg.LogStepProcess(startTime, "Step 3: Fetch DB Started")
	fetchService := fetch_db.NewFetchService(s.cfg, &dbRepository)
	err = fetchService.StartFetchDB(ctx, &objectProcessList)
	if err != nil {
		return err
	}
	pkg.LogStepProcess(startTime, "Step 3: Fetch DB Completed ")

	//Start Database Service & Process
	pkg.LogStepProcess(startTime, "Step 4: Import DB Started ")
	dbService := insert_data.NewDBService(s.cfg, &dbRepository, tempTarget)
	err = dbService.StartDBProcess(ctx, &objectProcessList)
	if err != nil {
		return err
	}
	pkg.LogStepProcess(startTime, "Step 4: Import DB Completed ")

	pkg.LogStepProcess(startTime, "Process Finished ")
	return nil
}

// Close closes server and related resources
func (s *Server) Close() {
	s.infra.Close()
	logger.AppLog(logger.InfoAppNewRecoverDataFinished, logger.InfoLevelLog)
}
