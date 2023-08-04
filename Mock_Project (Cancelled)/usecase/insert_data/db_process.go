package insert_data

import (
	"Mock_Project/infrastructure/repository"
	"Mock_Project/model"
	"Mock_Project/usecase/read_data"
	"context"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Server struct {
	config       *model.Server
	dbRepository repository.IDBRepository
	wg           *sync.WaitGroup
	tempPath     string
	err          chan error
}

func NewDBService(cfg *model.Server, dbRepository *repository.IDBRepository, path string) IDB {
	return &Server{
		config:       cfg,
		dbRepository: *dbRepository,
		wg:           new(sync.WaitGroup),
		tempPath:     path,
		err:          make(chan error, 1),
	}
}

func (s Server) StartDBProcess(ctx context.Context, collection *[]model.ObjectProcess) error {
	if len(*collection) == 0 {
		return fmt.Errorf("value is Empty")
	}
	prepareValue := s.prepareDBList(*collection)

	////Create Temp Folder
	fileService := read_data.NewService(s.tempPath, "")
	err := fileService.CreateParentFolder()
	if err != nil {
		return err
	}

	count := 0
	for _, collect := range prepareValue {
		//Detect Error in Goroutine
		err := s.breakError()
		if err != nil {
			return err
		}

		s.wg.Add(1)
		count += len(collect.Value)
		go s.processExportImport(ctx, collect)
	}
	s.wg.Wait()

	//Remove All Temp Folder
	defer func(name string) {
		_ = os.RemoveAll(name)
	}(s.tempPath)

	fmt.Println("Total Record ===> ", count)
	return nil
}

func (s Server) processExportImport(ctx context.Context, collect model.ObjectProcess) {
	defer s.wg.Done()
	file := s.tempPath + model.StrokeCharacter + collect.TableName

	//Get All Record from Table
	err := s.dbRepository.ExportDataFiles(file, ctx, collect)
	if err != nil {
		fmt.Println("Export Error ==>", err)
		s.err <- err
		return
	}

	//Add New Records to Temp Files
	var record []string
	for _, value := range collect.Value {
		str := convertObjectToString(value)
		record = append(record, str)
	}

	//Convert Object to Strings
	fileService := read_data.NewService(s.tempPath, collect.TableName)
	err = fileService.InsertCurrentFiles(&record)
	if err != nil {
		fmt.Println("Insert New Data Error ==>", err)
		s.err <- err
		return
	}

	//Truncate Remote all Current Data
	err = s.dbRepository.ClearTopic(ctx, collect)
	if err != nil {
		fmt.Println("Truncate Error ==>", err)
		s.err <- err
		return
	}

	//Import New Value to Table
	err = s.dbRepository.ImportDataFiles(file, ctx, collect)
	if err != nil {
		fmt.Println("Import Error ==>", err)
		s.err <- err
		return
	}
}

func (s Server) breakError() error {
	select {
	case err := <-s.err:
		return err
	default:
	}
	return nil
}

//func (s Server) processInsertData(ctx context.Context, collect model.ObjectProcess) {
//	err := s.dbRepository.InsertData(ctx, collect, []interface{}{})
//	if err != nil {
//		fmt.Println(err)
//	}
//	s.wg.Done()
//}

func (s Server) prepareDBList(collection []model.ObjectProcess) []model.ObjectProcess {
	for i := 0; i < len(collection); i++ {
		s.wg.Add(1)
		go s.sortItems(&collection[i])
	}
	s.wg.Wait()
	for i := 0; i < len(collection); i++ {
		for j := 0; j < len(collection[i].Value); j++ {
			curDate := formatDate(time.Now())
			curTime := formatTime(time.Now())
			collection[i].Value[j].TKSERIALNUMBER = 1
			collection[i].Value[j].TKZXD = curDate
			collection[i].Value[j].TKTIM = curTime
			if j > 0 {
				if collection[i].Value[j].QCD == collection[i].Value[j-1].QCD && curDate == collection[i].Value[j-1].TKZXD && curTime == collection[i].Value[j-1].TKTIM {
					collection[i].Value[j].TKSERIALNUMBER = collection[i].Value[j-1].TKSERIALNUMBER + 1
				}
			}
		}
	}
	return collection
}

func (s Server) sortItems(collect *model.ObjectProcess) {
	defer s.wg.Done()
	compare := func(i, j int) bool {
		return compareObject(i, j, collect.Value)
	}
	sort.Slice(collect.Value, compare)
}

func compareObject(i, j int, value []model.TargetObject) bool {
	if value[i].QCD != value[j].QCD {
		return strings.Compare(value[i].QCD, value[j].QCD) < 0
	}
	return strings.Compare(value[i].TIME, value[j].TIME) < 0
}

func formatTime(time time.Time) string {
	return time.Format(model.TimeFormatWithMicrosecond)
}

func formatDate(time time.Time) string {
	return time.Format(model.DateFormatWithStroke)
}

func convertObjectToString(object model.TargetObject) string {
	var temp []string
	val := reflect.ValueOf(&object).Elem()
	//typ := reflect.TypeOf(&object).Elem()
	for i := 0; i < val.NumField(); i++ {
		value := val.Field(i).Interface()
		switch v := value.(type) {
		case string:
			temp = append(temp, v)
		case int:
			temp = append(temp, strconv.Itoa(v))
		}
	}
	return strings.Join(temp, model.CommaCharacter)
}
