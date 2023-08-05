package db

import (
	"Mock_Project/model"
	"Mock_Project/pkg/logger"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type mysql struct {
	config  *model.Server
	clients map[string]*sql.DB
}

// NewDBHandler constructor
func NewDBHandler(cfg *model.Server) (IDBHandler, error) {
	return &mysql{
		config:  cfg,
		clients: make(map[string]*sql.DB),
	}, nil
}

func (c mysql) InitConnection(config *model.Server, endpoint, dbName string) error {
	key := endpoint + model.StrokeCharacter + dbName
	_, isExists := c.clients[key]
	if isExists {
		return nil
	}

	connectInfo := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s", config.User, config.Password, endpoint, config.Port, dbName,
	)

	client, err := sql.Open(config.DriverName, connectInfo)
	if err != nil {
		logger.AppLogDatabase(logger.ErrorCodeInitDatabaseConnectionFail, logger.ErrorLevelLog, endpoint, dbName)
		return err
	}

	if pingErr := client.Ping(); pingErr != nil {
		logger.AppLogDatabase(logger.ErrorCodeInitDatabaseConnectionFail, logger.ErrorLevelLog, endpoint, dbName)
		return pingErr
	}
	client.SetMaxOpenConns(config.MaxOpenConnection)
	client.SetMaxIdleConns(config.MaxIdleConnection)
	c.clients[key] = client

	return nil
}

func (c mysql) Exec(ctx context.Context, endpoint, dbName, sql string, args []interface{}) error {
	currentClient, isExists := c.clients[endpoint+model.StrokeCharacter+dbName]
	if !isExists {
		return fmt.Errorf("database connection not found")
	}
	tx, err := currentClient.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, sql, args...)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("transaction failed ===> %s", err)

	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit failed")
	}
	return nil
}

func (c mysql) CloseAllDb() error {
	for _, db := range c.clients {
		if err := db.Close(); err != nil {
			return err
		}
	}
	return nil
}
