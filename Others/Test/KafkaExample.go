package main

import (
	"Test/groupKafka"
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go groupKafka.Consumer()
	time.Sleep(1 * time.Second)
	go groupKafka.Producer()
	wg.Wait()
}
