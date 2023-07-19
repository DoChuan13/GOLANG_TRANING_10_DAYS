package CondExample

import (
    "fmt"
    "time"
)

func producer(out chan<- int) {
	for i := 0; i < 5; i++ {
		fmt.Println("producer, production: ", i)
		out <- i
	}
	close(out)
}

func consumer(in <-chan int) {
	for num := range in {
		fmt.Println("---consumer, consumption: ", num)
	}
}


func WithoutCond() {
	ch := make(chan int)
	go producer(ch)
	go consumer(ch)
	time.Sleep(5 * time.Second)
}
