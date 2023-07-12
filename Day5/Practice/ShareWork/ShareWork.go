package ShareWork

import (
	"Day5/Practice/WorkerPoolSelect/model"
	"fmt"
	"sync"
)

var allWorkers = 8
var workerOn = 5
var next = 5
var allJobs = 19
var jobList = make(chan int, allJobs)
var workers []model.Worker
var workGroup = make(chan model.Worker, workerOn)

func generateWorker() {
	for i := 0; i < allWorkers; i++ {
		var worker model.Worker = model.Person{Id: i}
		workers = append(workers, worker)
	}
}

func generateWorkGroup() {
	for i := 0; i < workerOn; i++ {
		workGroup <- workers[i]
	}
}

func nextWorkGroup() {
	fmt.Println("Next Worker start from ===>", next)
	workGroup <- workers[next]
	if next == allWorkers-1 {
		next = 0
	} else {
		next++
	}
}

func generateWorksChannel(wg *sync.WaitGroup) {
	for i := 1; i <= allJobs; i++ {
		jobList <- i
	}
	close(jobList)
	wg.Done()
}

func shareWorksForWorkers(wg *sync.WaitGroup) {
	isContinue := true
	var chWg sync.WaitGroup
	for i := 0; i < cap(workGroup); i++ {
		select {
		case v, ok := <-jobList:
			if ok {
				worker := <-workGroup
				chWg.Add(1)
				go worker.Worker(v, &chWg)
				nextWorkGroup()
			} else {
				isContinue = false
			}
		}
		if !isContinue {
			break
		}
		if i == cap(workGroup)-1 {
			i = -1
		}
	}
	chWg.Wait()
	wg.Done()
}

func MoreWorkerPool() {
	generateWorker()
	generateWorkGroup()
	var wg sync.WaitGroup
	wg.Add(2)

	go generateWorksChannel(&wg)
	go shareWorksForWorkers(&wg)

	wg.Wait()
}
