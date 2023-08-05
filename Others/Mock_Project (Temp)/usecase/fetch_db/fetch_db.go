package fetch_db

import (
	"Mock_Project/infrastructure/repository"
	"Mock_Project/model"
	"context"
)

type Server struct {
	config       *model.Server
	dbRepository repository.IDBRepository
}

func NewFetchService(cfg *model.Server, dbRepository *repository.IDBRepository) IFetchDB {
	return &Server{
		config:       cfg,
		dbRepository: *dbRepository,
	}
}

func (s Server) StartFetchDB(ctx context.Context, collection *model.ConsumerObject) error {
	err := s.dbRepository.InitConnection(s.config, s.config.Endpoint, s.config.DBName)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	err = s.dbRepository.CreateNewTable(ctx, *collection)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	return nil
}
