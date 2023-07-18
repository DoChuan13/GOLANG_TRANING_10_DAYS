package WaitGroup

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int) {
	fmt.Printf("Worker %d starting\n", id)

	time.Sleep(time.Second)
	fmt.Printf("Worker %d done\n", id)
}

func WGExample() {
	var wg sync.WaitGroup //Declaration WG
	for i := 1; i <= 5; i++ {
		wg.Add(1) //Add goroutine in WG, increase counter
		i := i
		go func() {
			defer wg.Done() //Confirm goroutine completed, decrease counter = wg.Add(-1)
			worker(i)
		}()
	}
	wg.Wait() //WG wait for until all goroutines finish, counter = 0
}
