package insert_data

import (
	"Mock_Project/infrastructure/repository"
	"Mock_Project/model"
	"Mock_Project/usecase/read_data"
	"context"
	"fmt"
	"sync"
)

type Server struct {
	config       *model.Server
	dbRepository repository.IDBRepository
	wg           *sync.WaitGroup
	err          chan error
}

func NewDBService(cfg *model.Server, dbRepository *repository.IDBRepository) IDB {
	return &Server{
		config:       cfg,
		dbRepository: *dbRepository,
		wg:           new(sync.WaitGroup),
		err:          make(chan error, 1),
	}
}

func (s Server) StartDBProcess(ctx context.Context, objectProcess *model.ConsumerObject) error {
	var err error = nil
	if len(objectProcess.Records) == 0 {
		return fmt.Errorf("value is Empty")
	}

	err = s.processExportImport(ctx, *objectProcess)
	if err != nil {
		return err
	}

	return nil
}

func (s Server) processExportImport(ctx context.Context, collect model.ConsumerObject) error {
	var err error = nil

	//1. Initial Connection
	err = s.dbRepository.InitConnection(s.config, s.config.Endpoint, s.config.DBName)
	if err != nil {
		return err
	}

	//2. GenerateTable And Get Current Record (Canceled Export File)
	err = s.dbRepository.GenerateTableAndExpFile(ctx, collect.TableName)
	if err != nil {
		fmt.Println("Generate Table Error ==>", err)
		s.err <- err
		return err
	}

	//3. Generate New Records to Temp Files
	fileService := read_data.NewService(s.config.LocalPath, collect.TableName)
	err = fileService.InsertCurrentFiles(&collect.Records)
	if err != nil {
		fmt.Println("Generate New Data Error ==>", err)
		s.err <- err
		return err
	}

	////4. Truncate Remote all Current Data (Canceled)
	//err = s.dbRepository.ClearData(ctx, collect.TableName)
	//if err != nil {
	//	fmt.Println("Truncate Error ==>", err)
	//	s.err <- err
	//	return err
	//}

	//5. Import New Value to Table
	err = s.dbRepository.ImportDataFiles(ctx, collect.TableName)
	if err != nil {
		fmt.Println("Import Error ==>", err)
		s.err <- err
		return err
	}
	return err
}
