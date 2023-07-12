package GoRoutines

import (
	"fmt"
	"time"
)

func GoRoutine() {
	w1 := worker{name: "A"}
	w2 := worker{name: "B"}
	w3 := worker{name: "C"}
	w4 := worker{name: "D"}
	w5 := worker{name: "E"}
	//List Worker
	workers := []worker{w1, w2, w3, w4, w5}

	//Report Worker
	var report [][]int

	//List All Work
	allWork := 20
	workList := make(chan int, allWork)
	for i := 1; i <= allWork; i++ {
		workList <- i
	}
	close(workList)

	//List Task for Worker
	sizeWork := allWork / len(workers)
	if allWork%len(workers) != 0 {
		sizeWork++
	}
	var taskList []chan int
	for i := 0; i < len(workers); i++ {
		taskList = append(taskList, make(chan int, sizeWork))
	}

	//Share Work for a Worker
	for i, j := 0, 0; i < cap(workList) && j < len(workers); i, j = i+1, j+1 {
		taskList[j] <- <-workList
		if j == len(workers)-1 {
			j = -1
		}
	}
	for i := 0; i < len(taskList); i++ {
		close(taskList[i])
	}

	//All Worker working
	for i := 0; i < len(workers); i++ {
		num := 0
		var taskGroup []int
		for task := range taskList[i] {
			go workers[i].Work1(task)
			taskGroup = append(taskGroup, task)
			num++
		}
		report = append(report, taskGroup)
	}
	time.Sleep(1 * time.Second)
	fmt.Println("Report for all worker", report)
}

type worker struct {
	name string
}

type work interface {
	Work(channel chan int)
	work1(channel int)
}

func (w worker) Work(channel chan int) {
	fmt.Printf("worker %s is working ===> %d\n", w.name, <-channel)
}
func (w worker) Work1(channel int) {
	fmt.Printf("worker %s is working ===> %d\n", w.name, channel)
}
