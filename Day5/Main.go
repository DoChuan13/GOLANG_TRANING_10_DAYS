package main

import (
	"Day5/Practice/ShareWork"
	"fmt"
)

func main() {
	fmt.Println("======================1. Example Goroutine=======================")
	//Phân chia công việc cho worker
	//GoRoutines.GoRoutine()
	fmt.Println("======================2. Example WaitGroup 1=======================")
	//Wait Group
	//Wait.Group()
	fmt.Println("======================2. Example WaitGroup 2=======================")
	//
	//Wait.Wait()
	fmt.Println("======================3. Example Simple Sync Channel=======================")
	//Đồng bộ hóa các Goroutines
	//Sync.SyncChannel()
	fmt.Println("======================4. Example Select Case=======================")
	//Ví dụ Select Case
	//Select.Select()
	fmt.Println("======================5. Example Select Worker=======================")
	//Chia đều công việc cho các worker
	//WorkerPoolSelect.MoreWorkerPoolExe()
	fmt.Println("======================5. Example Select Worker 2=======================")
	//Giới hạn số worker có thể thực hiện công việc
	ShareWork.MoreWorkerPool()
}
