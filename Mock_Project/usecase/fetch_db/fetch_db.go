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
}

func NewFetchService(cfg *model.Server, dbRepository *repository.IDBRepository) IFetchDB {
	return &Server{
		config:       cfg,
		dbRepository: *dbRepository,
		wg:           new(sync.WaitGroup),
	}
}

func (s Server) StartFetchDB(ctx context.Context, collection *[]model.ObjectProcess) error {
	for _, collect := range *collection {
		err := s.initConnection(collect)
		if err != nil {
			return err
		}
		err = s.dbRepository.CreateNewTable(ctx, collect)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s Server) initConnection(collect model.ObjectProcess) error {
	err := s.dbRepository.InitConnection(s.config, collect.EndPoint, collect.DBName)
	if err != nil {
		return err
	}
	return nil
}
