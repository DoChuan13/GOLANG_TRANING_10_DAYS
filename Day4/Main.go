package main

import (
	"fmt"
	"time"
)

func numbers() {
	for i := 1; i <= 5; i++ {
		time.Sleep(250 * time.Millisecond)
		fmt.Printf("%d ", i)
	}
}
func alphabets() {
	for i := 'a'; i <= 'e'; i++ {
		time.Sleep(400 * time.Millisecond)
		fmt.Printf("%c ", i)
	}
}
func main() {
	fmt.Println("========================1 Goroutines========================")
	go numbers()
	go alphabets()
	fmt.Println("After All Go")
	time.Sleep(3000 * time.Millisecond)
	fmt.Println("main terminated")

	fmt.Println("========================2. Channel========================")
	var cha chan int
	fmt.Println(cha)
	myChan := make(chan int)

	go func() {
		send := 1
		fmt.Println("Send value to Channel", send)
		myChan <- send
	}()
	fmt.Println("Receive Value from Channel", <-myChan)

	myChanSy := make(chan int)

	go func() {
		for i := 1; i <= 5; i++ {
			myChanSy <- i
			time.Sleep(400 * time.Millisecond)
		}
	}()

	for i := 1; i <= 5; i++ {
		fmt.Println(<-myChanSy)
	}
	fmt.Println("========================2.1. Channel========================")

	demo := make(chan int)
	go getChannel(demo)

	for channel := range demo {
		fmt.Println("Channel value =>>>", channel)
	}

	myChan2 := make(chan int)

	go func() {
		for i := 1; i <= 10; i++ {
			myChan2 <- i
		}
		close(myChan2)
	}()

	for {
		value, isAlive := <-myChan2
		if !isAlive {
			fmt.Printf("Value: %d. Channel has been closed.\n", value)
			break
		}
		fmt.Printf("Value: %d\n", value)
	}
	time.Sleep(time.Second)
	buffCh := make(chan int, 3)
	buffCh <- 1
	buffCh <- 2
	fmt.Println("Capacity of Buffer Channel===>", cap(buffCh))
	fmt.Println("Length of Buffer Channel===>", len(buffCh))
}

// Function for Goroutines
func getChannel(ch chan int) {
	for i := 0; i < 5; i++ {
		ch <- i
		//fmt.Println("Set Channel value success")
	}
	close(ch)
}
