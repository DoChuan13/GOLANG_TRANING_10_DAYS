package Mutex

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type Container struct {
	mu       sync.Mutex
	counters map[string]int
}

func (c *Container) inc(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counters[name]++
}

func MutexExample() {
	c := Container{
		counters: map[string]int{"a": 0, "b": 0},
	}
	var wg sync.WaitGroup

	doIncrement := func(name string, n int) {
		for i := 0; i < n; i++ {
			c.inc(name)
		}
		wg.Done()
	}

	wg.Add(3)
	go doIncrement("a", 10000)
	go doIncrement("a", 10000)
	go doIncrement("b", 10000)

	wg.Wait()
	fmt.Println(c.counters)
}

var mu sync.Mutex

// Khai báo biến count được truy cập bởi tất cả các routine
var count = 0

func process(n int, wg *sync.WaitGroup) {
	for i := 0; i < 10; i++ {
		time.Sleep(time.Duration(rand.Int31n(2)) * time.Second)
		mu.Lock() //Bắt đầu khoá ở đây
		temp := count
		temp++
		time.Sleep(time.Duration(rand.Int31n(2)) * time.Second)
		count = temp
		mu.Unlock() // Mở khoá
	}
	fmt.Println("Count after i="+strconv.Itoa(n)+" Count:", strconv.Itoa(count))
	wg.Done()
}

func MutexExample2() {
	var wg sync.WaitGroup
	for i := 1; i < 4; i++ {
		wg.Add(1)
		go process(i, &wg)
	}
	wg.Wait() //Tạm dừng để đợi cho tất cả routine hoàn thành
	fmt.Println("Final Count:", count)
}
