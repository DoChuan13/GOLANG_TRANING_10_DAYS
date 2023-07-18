package main

import (
	"TestConnect/config"
	"TestConnect/database"
	"TestConnect/model"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup
var filePath = "model/save.csv"
var list = config.ReadFile(filePath)
var userListChan = make(chan model.User, len(list))
var ready = make(chan bool, 1)
var cancelRoutine = make(chan bool, 1)
var workList []model.User
var isContinue = true

func main() {
	db := *config.Connection()
	defer config.Close(&db)
	database.CreateTables(&db)
	//ReadWriteFile()

	generateSaveList(&userListChan, &list)

	wg.Add(1)
	go preparePostDBList()
	ticker := time.NewTicker(3 * time.Second)
	for isContinue {
		//fmt.Println("In For loop, length==> , remain ", len(userListChan))
		if !isContinue {
			fmt.Println("Stop Loop and Ticker")
			ticker.Stop()
			return
		}
		select {
		case <-ready:
			cancelRoutine <- true
			fmt.Println("WorkList Case 1 (max 10) ===>", len(workList))
			for i := 0; i < len(workList); i++ {
				curUser := workList[i]
				database.SaveDataDB(&db, &curUser)
			}
			if len(userListChan) == 0 {
				fmt.Println("Switch continue = false (max10)")
				isContinue = false
			}
			//fmt.Println(" ===> Is Continue WorkList Case 1 (max 10), remain", len(userListChan))
			workList = []model.User{}
			if isContinue {
				<-cancelRoutine
				wg.Add(1)
				go preparePostDBList()
			}
		case <-ticker.C:
			cancelRoutine <- true
			fmt.Println("WorkList Case 2 (ticker) ===>", len(workList))
			for i := 0; i < len(workList); i++ {
				curUser := workList[i]
				database.SaveDataDB(&db, &curUser)
			}
			if len(userListChan) == 0 {
				fmt.Println("Switch continue = false(ticker)")
				isContinue = false
			}
			//fmt.Println(" ===>Is Continue WorkList Case 2 (ticker), remain", len(userListChan))
			workList = []model.User{}
			if isContinue {
				<-cancelRoutine
				wg.Add(1)
				go preparePostDBList()
			}
		}
	}
	wg.Wait()
	fmt.Println("Program finished!!!")
}

func generateSaveList(channel *chan model.User, userList *[]string) {
	for _, value := range *userList {
		var user model.User
		_ = json.Unmarshal([]byte(value), &user)
		*channel <- user
	}
	close(*channel)
}

func preparePostDBList() {
	//fmt.Printf("==> Prepare  ==>%d, remain %d\n", len(workList), len(userListChan))
	for true {
		select {
		case <-cancelRoutine:
			wg.Done()
			return
		default:
			value, ok := <-userListChan
			if ok {
				workList = append(workList, value)
				time.Sleep(200 * time.Millisecond)
				if len(workList) == 10 {
					fmt.Println("Prepare Ready")
					ready <- true
					wg.Done()
					return
				}
			} else {
				ready <- true
				wg.Done()
				return
			}
		}
	}
}

func ReadWriteFile() {
	filePath := "model/save.csv"
	var user model.User
	cur := config.ReadFile(filePath)
	id := 0
	if len(cur) != 0 {
		_ = json.Unmarshal([]byte(cur[len(cur)-1]), &user)
		id = user.Id
	}
	for {
		id++
		user.Id = id
		fmt.Print("Input Name: ")
		user.Name = config.InputString()
		result, _ := json.Marshal(user)
		fmt.Println(string(result))
		file := config.OpenFile(filePath)
		config.WriteFile(file, []string{string(result)})
	}
}
