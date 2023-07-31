package repository

import (
	"Mock_Project/infrastructure"
	"Mock_Project/infrastructure/db"
	"Mock_Project/model"
	"context"
	"fmt"
	"reflect"
	"strconv"
)

const baseInsertIntoTable = "insert into %s %s values %s;"

type dbRepository struct {
	config model.Server
	db     db.IDBHandler
}

func NewDBRepository(infra *infrastructure.Infra, cfg model.Server) IDBRepository {
	return &dbRepository{
		config: cfg,
		db:     infra.DBHandler,
	}
}

func (r dbRepository) ImportData(
	ctx context.Context, object model.ObjectProcess, tableName, endpoint, dbName string, args []interface{},
) error {
	if len(object.Value) == 0 {
		return fmt.Errorf("value is Empty")
	}
	columns := ConvertColumns(object.Value[0])
	values := ConvertValues(object.Value)
	query := fmt.Sprintf(baseInsertIntoTable, tableName, columns, values)
	err := r.db.Exec(ctx, endpoint, dbName, query, []interface{}{})
	if err != nil {
		return err
	}
	return nil
}

func ConvertColumns(object model.TargetObject) string {
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

func ConvertValues(collect []model.TargetObject) interface{} {
	result := ""
	for i := 0; i < len(collect); i++ {
		result += model.OpenRoundBracket
		val := reflect.ValueOf(collect[i])
		for j := 0; j < val.NumField(); j++ {
			temp := val.Field(j).Interface()
			switch v := temp.(type) {
			case string:
				result += v
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
}

func (r dbRepository) Close() error {
	return r.db.Close()
}
