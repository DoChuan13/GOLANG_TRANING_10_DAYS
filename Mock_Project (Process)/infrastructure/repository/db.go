package repository

import (
	"Mock_Project/infrastructure"
	"Mock_Project/infrastructure/db"
	"Mock_Project/model"
	"context"
	"fmt"
)

const (
	//baseInsertIntoTable = "insert into %s.%s %s values %s;"
	baseCreateTableAndExportFile = "call GenerateTableAndGetRecord('%s','%s');"
	baseLoadImportFiles          = "load data infile '%s' into table %s fields terminated by ',' lines terminated by '\n';"
	baseQueryClearData           = "truncate table %s;"
)

type dbRepository struct {
	config *model.Server
	db     db.IDBHandler
}

func NewDBRepository(infra *infrastructure.Infra, cfg *model.Server) IDBRepository {
	return &dbRepository{
		config: cfg,
		db:     infra.DBHandler,
	}
}

func (r dbRepository) InitConnection(config *model.Server, endpoint, dbName string) error {
	return r.db.InitConnection(config, endpoint, dbName)
}

func (r dbRepository) GenerateTableAndExpFile(ctx context.Context, tableName string) error {
	file := r.config.SqlPath + tableName
	query := fmt.Sprintf(baseCreateTableAndExportFile, tableName, file)
	err := r.db.Exec(ctx, r.config.Endpoint, r.config.DBName, query, []interface{}{})
	if err != nil {
		return err
	}
	return nil
}

func (r dbRepository) ImportDataFiles(ctx context.Context, tableName string) error {
	file := r.config.SqlPath + tableName
	query := fmt.Sprintf(baseLoadImportFiles, file, tableName)
	err := r.db.Exec(ctx, r.config.Endpoint, r.config.DBName, query, []interface{}{})
	if err != nil {
		return err
	}
	return nil
}

func (r dbRepository) ClearData(ctx context.Context, tableName string) error {
	query := fmt.Sprintf(baseQueryClearData, tableName)
	err := r.db.Exec(ctx, r.config.Endpoint, r.config.DBName, query, []interface{}{})
	if err != nil {
		return err
	}
	return nil
}

func (r dbRepository) CloseAllDb() error {
	return r.db.CloseAllDb()
}
