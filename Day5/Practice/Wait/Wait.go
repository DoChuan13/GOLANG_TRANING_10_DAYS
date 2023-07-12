package Wait

import (
	"fmt"
	"sync"
	"time"
)

func printNumber(wg *sync.WaitGroup) {
	for i := 0; i < 10; i++ {
		fmt.Println("Print Number =>>>", i)
		time.Sleep(time.Millisecond * 400)
	}
	wg.Done() //5.1 Thông báo hoàn tất
}

func printCharacter(wg *sync.WaitGroup) {
	hello := "Hello World"
	for _, value := range hello {
		fmt.Println("Print Character =>>>", string(value))
		time.Sleep(time.Millisecond * 200)
	}
	wg.Done() //5.2 Thông báo hoàn tất
}

func Wait() {
	var wg sync.WaitGroup //1. Khởi tạo WG
	wg.Add(2)             //2. Thiết lập chờ 2 goroutines

	go printNumber(&wg)    //3.1Gửi WG cho goroutines
	go printCharacter(&wg) //3.2Gửi WG cho goroutines

	wg.Wait() //4. Thực hiện chờ
}
