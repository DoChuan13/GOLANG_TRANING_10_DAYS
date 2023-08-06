package insert_data

import (
	"Mock_Project/infrastructure/repository"
	"Mock_Project/model"
	"Mock_Project/usecase/read_data"
	"context"
	"fmt"
	"os"
	"sync"
)

type Server struct {
	config       *model.Server
	dbRepository repository.IDBRepository
	wg           *sync.WaitGroup
	tempPath     string
	err          chan error
}

func NewDBService(cfg *model.Server, dbRepository *repository.IDBRepository, path string) IDB {
	return &Server{
		config:       cfg,
		dbRepository: *dbRepository,
		wg:           new(sync.WaitGroup),
		tempPath:     path,
		err:          make(chan error, 1),
	}
}

func (s Server) StartDBProcess(ctx context.Context, objectProcess *model.ConsumerObject) error {
	if len(objectProcess.Records) == 0 {
		return fmt.Errorf("value is Empty")
	}

	////Create Temp Folder
	fileService := read_data.NewService(s.tempPath, "")
	err := fileService.CreateParentFolder()
	if err != nil {
		return err
	}

	s.processExportImport(ctx, *objectProcess)

	//Remove Temp Folder
	err = os.Remove(s.tempPath + model.StrokeCharacter + objectProcess.TableName)
	if err != nil {
		return err
	}

	return nil
}

func (s Server) processExportImport(ctx context.Context, collect model.ConsumerObject) {
	file := s.tempPath + model.StrokeCharacter + collect.TableName
	var err error = nil

	//1. Initial Connection
	err = s.dbRepository.InitConnection(s.config, s.config.Endpoint, s.config.DBName)
	if err != nil {
		return
	}

	//2. GenerateTable And Get Current Record
	err = s.dbRepository.GenerateTableAndExpFile(file, ctx, collect)
	if err != nil {
		fmt.Println("Generate Table Error ==>", err)
		s.err <- err
		return
	}

	//3. Add New Records to Temp Files
	fileService := read_data.NewService(s.tempPath, collect.TableName)
	err = fileService.InsertCurrentFiles(&collect.Records)
	if err != nil {
		fmt.Println("Insert New Data Error ==>", err)
		s.err <- err
		return
	}

	//4. Truncate Remote all Current Data
	err = s.dbRepository.ClearData(ctx, collect)
	if err != nil {
		fmt.Println("Truncate Error ==>", err)
		s.err <- err
		return
	}

	//5. Import New Value to Table
	err = s.dbRepository.ImportDataFiles(file, ctx, collect)
	if err != nil {
		fmt.Println("Import Error ==>", err)
		s.err <- err
		return
	}
}
