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
	baseCreateTable     = "create table if not exists %s.%s (%s)"
	baseLoadImportFiles = "load data infile '%s' into table %s.%s fields terminated by ',' lines terminated by '\n';"
	baseLoadExportFiles = "select * into outfile '%s' fields terminated by ',' lines terminated by '\n' from %s.%s;"
	baseQueryClearData  = "truncate table %s.%s;"
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

func (r dbRepository) CreateNewTable(ctx context.Context, object model.ConsumerObject) error {
	fields := generateFields()
	query := fmt.Sprintf(baseCreateTable, r.config.DBName, object.TableName, fields)
	err := r.db.Exec(ctx, r.config.Endpoint, r.config.DBName, query, []interface{}{})
	if err != nil {
		return err
	}
	return nil
}

func (r dbRepository) ImportDataFiles(file string, ctx context.Context, object model.ConsumerObject) error {
	query := fmt.Sprintf(baseLoadImportFiles, file, r.config.DBName, object.TableName)
	err := r.db.Exec(ctx, r.config.Endpoint, r.config.DBName, query, []interface{}{})
	if err != nil {
		return err
	}
	return nil
}

func (r dbRepository) ExportDataFiles(file string, ctx context.Context, object model.ConsumerObject) error {
	query := fmt.Sprintf(baseLoadExportFiles, file, r.config.DBName, object.TableName)
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

func generateFields() string {
	result := `QCD            varchar(42) not null,
            TIME           varchar(5)  not null,
            TKZXD          date        not null,
            TKTIM          varchar(15) not null,
            TKSERIALNUMBER int         not null,
            TKQKBN         varchar(2)  not null,
            SNDC           varchar(5)  not null,
            ZXD            varchar(10) null,
            NOPV           varchar(24) null,
            ZSS            varchar(4)  null,
            DPS            varchar(4)  null,
            DPS_T          varchar(5)  null,
            DPP            varchar(24) null,
            DPP_T          varchar(5)  null,
            DPR            varchar(24) null,
            DPR_T          varchar(5)  null,
            XV             varchar(24) null,
            XV_T           varchar(5)  null,
            NOQ            varchar(24) null,
            QAS            varchar(4)  null,
            QAP            varchar(24) null,
            QAP_T          varchar(5)  null,
            QAR            varchar(24) null,
            QAR_T          varchar(5)  null,
            AV             varchar(24) null,
            AV_T           varchar(5)  null,
            QBS            varchar(4)  null,
            QBP            varchar(24) null,
            QBP_T          varchar(5)  null,
            QBR            varchar(24) null,
            QBR_T          varchar(5)  null,
            BV             varchar(24) null,
            BV_T           varchar(5)  null,
            J_SNDC         varchar(5)  null,
            XJ             varchar(24) null,
            XJ_T           varchar(5)  null,
            DPCF           varchar(5)  null,
            JQCS           varchar(4)  null,
            JQJS           varchar(4)  null,
            DPCY           varchar(24) null,
            DPCY_T         varchar(5)  null,
            TLNM           varchar(24) null,
            XV_S           varchar(8)  null,
            AV_S           varchar(8)  null,
            BV_S           varchar(8)  null,
            DPP_S          varchar(8)  null,
            J_ZXD          varchar(10) null,
            MIDP           varchar(24) null,
            MIDP_T         varchar(5)  null,
            DYWP           varchar(24) null,
            DYWP_T         varchar(5)  null,
            DYWR           varchar(24) null,
            VWAP           varchar(24) null,
            VWAP_T         varchar(5)  null,
            ABV            varchar(24) null,
            ABV_T          varchar(5)  null,
            AAV            varchar(24) null,
            AAV_T          varchar(5)  null,
            QOV            varchar(24) null,
            QOV_T          varchar(5)  null,
            QUV            varchar(24) null,
            QUV_T          varchar(5)  null,
            INAV           varchar(24) null,
            INAV_T         varchar(5)  null,
            IYRP           varchar(24) null,
            IYRP_T         varchar(5)  null,
            IQRP           varchar(24) null,
            IQRP_T         varchar(5)  null,
            QACY           varchar(24) null,
            QACY_T         varchar(5)  null,
            QBCY           varchar(24) null,
            QBCY_T         varchar(5)  null,
            TSTF           varchar(4)  null,
            NO             varchar(24) null,` +
		"`NOT`          varchar(24) null," +
		`1KF          varchar(4)  null,
        2KF          varchar(4)  null,
        DKF            varchar(4)  null,
        DSYRP          varchar(24) null,
        DSYWP          varchar(24) null,
        primary key (QCD, TKZXD, TKTIM, TKSERIALNUMBER)`
	return result
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
