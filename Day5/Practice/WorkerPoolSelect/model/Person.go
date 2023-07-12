package model

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Person struct {
	Id int
}

func (p Person) Worker(job int, wg *sync.WaitGroup) {
	var t int = int(rand.Intn(999))
	time.Sleep(time.Duration(t))
	fmt.Printf("Worker %d is working job %d with time %d\n", p.Id, job, t)
	wg.Done()
}
