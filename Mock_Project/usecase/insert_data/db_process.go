package insert_data

import (
	"Mock_Project/infrastructure/repository"
	"Mock_Project/model"
	"context"
	"fmt"
	"sync"
)

type Server struct {
	config       *model.Server
	dbRepository repository.IDBRepository
	wg           *sync.WaitGroup
}

func NewService(cfg *model.Server, dbRepository repository.IDBRepository) IDB {
	return &Server{
		config:       cfg,
		dbRepository: dbRepository,
		wg:           new(sync.WaitGroup),
	}
}

func (s Server) StartDBProcess(ctx context.Context, collection *[]model.ObjectProcess) error {
	if len(*collection) == 0 {
		return fmt.Errorf("value is Empty")
	}
	for _, collect := range *collection {
		s.wg.Add(1)
		go s.processInsertData(ctx, collect)
	}
	s.wg.Wait()
	return nil
}

func (s Server) processInsertData(ctx context.Context, collect model.ObjectProcess) {
	err := s.dbRepository.ImportData(
		ctx, collect, collect.TableName, collect.EndPoint, collect.DBName, []interface{}{},
	)
	if err != nil {
		fmt.Println("Error")
	}
	s.wg.Done()
}
