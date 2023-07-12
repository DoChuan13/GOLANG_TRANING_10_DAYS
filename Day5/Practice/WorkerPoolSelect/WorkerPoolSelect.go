package WorkerPoolSelect

import (
	"Day5/Practice/WorkerPoolSelect/model"
	"sync"
)

var allWorkers = 5
var allJobs = 20
var jobList = make(chan int, allJobs)
var workers []model.Worker

func generateWorker() {
	for i := 0; i < allWorkers; i++ {
		var worker model.Worker = model.Person{Id: i}
		workers = append(workers, worker)
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
	for i := 0; i < allWorkers; i++ {
		select {
		case v, ok := <-jobList:
			if ok {
				chWg.Add(1)
				go workers[i].Worker(v, &chWg)
			} else {
				isContinue = false
			}
		}
		if !isContinue {
			break
		}
		if i == allWorkers-1 {
			i = -1
		}
	}
	chWg.Wait()
	wg.Done()
}

func WorkerPoolExe() {
	generateWorker()
	var wg sync.WaitGroup
	wg.Add(2)

	go generateWorksChannel(&wg)
	go shareWorksForWorkers(&wg)

	wg.Wait()
}
