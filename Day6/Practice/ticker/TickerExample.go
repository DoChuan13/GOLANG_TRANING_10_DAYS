package ticker

import (
	"fmt"
	"time"
)

func Ticker() {
	ticker := time.NewTicker(500 * time.Millisecond) // =>Khởi tạo ticker thiết lập khoảng thời gian
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C: // =>Ticker thực thi theo khoảng thời gian định trước
				fmt.Println("Tick at", t)
			}
		}
	}()
	time.Sleep(2600 * time.Millisecond)
	ticker.Stop()
	done <- true
	fmt.Println("Ticker stopped")
}
