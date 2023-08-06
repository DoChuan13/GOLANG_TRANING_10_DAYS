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
	file := r.config.DockerPath + object.TableName
	query := fmt.Sprintf(baseCreateTableAndExportFile, object.TableName, file)
	err := r.db.Exec(ctx, r.config.Endpoint, r.config.DBName, query, []interface{}{})
	if err != nil {
		return err
	}
	return nil
}

func (r dbRepository) ImportDataFiles(ctx context.Context, object model.ConsumerObject) error {
	file := r.config.LocalPath + object.TableName
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

/*func (r dbRepository) InsertData(
	ctx context.Context, objectProcess model.ConsumerObject, args []interface{},
) error {
	if len(objectProcess.Records) == 0 {
		return fmt.Errorf("value is Empty")
	}
	columns := convertColumns(objectProcess.Value[0])
	values := convertValues(objectProcess.Value)
	query := fmt.Sprintf(baseInsertIntoTable, objectProcess.DBName, objectProcess.TableName, columns, values)
	err := r.db.Exec(ctx, objectProcess.EndPoint, objectProcess.DBName, query, args)
	if err != nil {
		return err
	}
	return nil
}

func convertColumns(object model.TargetObject) string {
	val := reflect.ValueOf(&object).Elem()
	typ := reflect.TypeOf(&object).Elem()
	columns := model.OpenRoundBracket
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		temp := field.Name
		if temp == "K1KF" {
			columns += "1KF"
		} else if temp == "K2KF" {
			columns += "2KF"
		} else if temp == "NOT" {
			columns += "`NOT`"
		} else {
			columns += temp
		}
		if i < val.NumField()-1 {
			columns += model.CommaCharacter
		}
	}
	columns += model.CloseRoundBracket
	return columns
}

func convertValues(collect []model.TargetObject) interface{} {
	result := ""
	for i := 0; i < len(collect); i++ {
		result += model.OpenRoundBracket
		val := reflect.ValueOf(collect[i])
		for j := 0; j < val.NumField(); j++ {
			temp := val.Field(j).Interface()
			switch v := temp.(type) {
			case string:
				result += model.ApostropheCharacter + v + model.ApostropheCharacter
			case int:
				result += strconv.Itoa(v)
			}

			if j < val.NumField()-1 {
				result += model.CommaCharacter
			}
		}
		result += model.CloseRoundBracket
		if i != len(collect)-1 {
			result += model.CommaCharacter + model.NewLineCharacter
		}
	}
	return result
}*/
