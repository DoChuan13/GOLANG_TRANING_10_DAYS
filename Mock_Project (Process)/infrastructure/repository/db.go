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
	baseLoadImportFiles          = "load data infile '%s' into table %s.%s fields terminated by ',' lines terminated by '\n';"
	baseQueryClearData           = "truncate table %s.%s;"
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

func (r dbRepository) GenerateTableAndExpFile(ctx context.Context, object model.ConsumerObject) error {
	file := r.config.SqlPath + object.TableName
	query := fmt.Sprintf(baseCreateTableAndExportFile, object.TableName, file)
	err := r.db.Exec(ctx, r.config.Endpoint, r.config.DBName, query, []interface{}{})
	if err != nil {
		return err
	}
	return nil
}

func (r dbRepository) ImportDataFiles(ctx context.Context, object model.ConsumerObject) error {
	file := r.config.SqlPath + object.TableName
	query := fmt.Sprintf(baseLoadImportFiles, file, r.config.DBName, object.TableName)
	err := r.db.Exec(ctx, r.config.Endpoint, r.config.DBName, query, []interface{}{})
	if err != nil {
		return err
	}
	return nil
}

func (r dbRepository) ClearData(ctx context.Context, object model.ConsumerObject) error {
	query := fmt.Sprintf(baseQueryClearData, r.config.DBName, object.TableName)
	err := r.db.Exec(ctx, r.config.Endpoint, r.config.DBName, query, []interface{}{})
	if err != nil {
		return err
	}
	return nil
}

func (r dbRepository) CloseAllDb() error {
	return r.db.CloseAllDb()
}
