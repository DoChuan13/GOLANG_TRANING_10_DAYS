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
	"sync"
	"time"
)

type Server struct {
	infra *infrastructure.Infra
	cfg   *model.Server
	cond  sync.Cond
	ltdDB chan bool
}

// New create new server instance
func New(infra *infrastructure.Infra, cfg *model.Server) *Server {
	return &Server{
		infra: infra,
		cfg:   cfg,
		cond:  sync.Cond{},
		ltdDB: make(chan bool, 400),
	}
}

var wg sync.WaitGroup

func (s *Server) Start(ctx context.Context) error {
	//Config
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	tempTarget := dir + "/file/temp"
	var mkdirPath = "file/collect"
	//var fileName = "A0001&0F"

	fmt.Print("Input File Name: ")
	var file string
	_, _ = fmt.Scanf("%s", &file)

	//file.FakeAllData()

	var startTime = time.Now()
	//Read File Process
	pkg.LogStepProcess(startTime, "Step 1: Read File Started")
	fileService := read_data.NewService(mkdirPath, file)
	_ = fileService.RemoveFolder(tempTarget)
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
	var consumerCh = make(chan model.ConsumerObject)
	var done = make(chan bool, 1)
	//var done = make(chan bool, 1)
	go func() {
		err = kafkaService.StartKafkaProcess(consumerCh, done, rows)
	}()

	//Start Fetch Database Service & Process
	//Start Database Service & Process
	fetchService := fetch_db.NewFetchService(s.cfg, &dbRepository)
	dbService := insert_data.NewDBService(s.cfg, &dbRepository, tempTarget)

	isContinue := true
	counter := 0
	s.cond.L = &sync.Mutex{}
	for isContinue {
		select {
		case objectProcess := <-consumerCh:
			s.ltdDB <- true
			counter += len(objectProcess.Records)
			//if counter%50000 == 0 {
			//	fmt.Println("Check point", counter)
			//}
			s.cond.L.Lock()
			if len(s.ltdDB) == cap(s.ltdDB) {
				s.cond.Wait()
			}
			wg.Add(1)
			go s.processDatabase(ctx, fetchService, &objectProcess, dbService)
			s.cond.L.Unlock()

		case val := <-done:
			fmt.Printf("Kafka Process Completed => %t", val)
			isContinue = false
			close(done)
		}
	}

	wg.Wait()
	close(consumerCh)

	pkg.LogStepProcess(startTime, "Step 3: Import DB Completed ")
	fmt.Println("Total Count ===>", counter)

	pkg.LogStepProcess(startTime, "Process Finished ")
	return nil
}

func (s *Server) processDatabase(
	ctx context.Context, fetchService fetch_db.IFetchDB, objectProcess *model.ConsumerObject, dbService insert_data.IDB,
) {
	defer wg.Done()
	err := fetchService.StartFetchDB(ctx, objectProcess)
	if err != nil {
		return
	}
	err = dbService.StartDBProcess(ctx, objectProcess)
	if err != nil {
		return
	}
	//fmt.Println("Saved Table =>", objectProcess.TableName)
	<-s.ltdDB
	s.cond.Signal()
}

// Close closes server and related resources
func (s *Server) Close() {
	s.infra.Close()
	logger.AppLog(logger.InfoAppNewRecoverDataFinished, logger.InfoLevelLog)
}
