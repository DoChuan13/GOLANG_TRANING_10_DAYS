package Wait

import (
	"fmt"
	"sync"
	"time"
)

func process(i int, wg *sync.WaitGroup) {
	fmt.Println("started Goroutine ", i)
	time.Sleep(2 * time.Second)
	fmt.Printf("Goroutine %d ended\n", i)
	wg.Done()
}

func Group() {
	no := 3
	var wg sync.WaitGroup
	for j := 0; j < 4; j++ {
		for i := 0; i < no; i++ {
			wg.Add(1)
			go process(i, &wg)
		}
	}
	wg.Wait()
	fmt.Println("All go routines finished executing")
}
