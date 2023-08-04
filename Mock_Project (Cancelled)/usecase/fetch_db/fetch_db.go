package fetch_db

import (
	"Mock_Project/infrastructure/repository"
	"Mock_Project/model"
	"context"
	"sync"
)

type Server struct {
	config       *model.Server
	dbRepository repository.IDBRepository
	wg           *sync.WaitGroup
	mu           *sync.Mutex
	err          chan error
}

func NewFetchService(cfg *model.Server, dbRepository *repository.IDBRepository) IFetchDB {
	return &Server{
		config:       cfg,
		dbRepository: *dbRepository,
		wg:           new(sync.WaitGroup),
		mu:           new(sync.Mutex),
		err:          make(chan error, 1),
	}
}

func (s Server) StartFetchDB(ctx context.Context, collection *[]model.ObjectProcess) error {
	for _, collect := range *collection {
		//Detect Error in Goroutine
		err := s.breakError()
		if err != nil {
			return err
		}
		s.wg.Add(1)
		go s.processGenerateTable(ctx, err, collect)
	}
	s.wg.Wait()
	return nil
}

func (s Server) initConnection(collect model.ObjectProcess) error {
	err := s.dbRepository.InitConnection(s.config, collect.EndPoint, collect.DBName)
	if err != nil {
		return err
	}
	return nil
}

func (s Server) processGenerateTable(ctx context.Context, err error, collect model.ObjectProcess) {
	defer s.wg.Done()
	s.mu.Lock()
	err = s.initConnection(collect)
	if err != nil {
		s.err <- err
		return
	}
	err = s.dbRepository.CreateNewTable(ctx, collect)
	if err != nil {
		s.err <- err
		return
	}
	s.mu.Unlock()
}

func (s Server) breakError() error {
	select {
	case err := <-s.err:
		return err
	default:
	}
	return nil
}
