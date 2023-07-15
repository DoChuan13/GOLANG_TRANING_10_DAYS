package panic

import "fmt"

func fullName(firstName *string, lastName *string) {
	if firstName == nil {
		panic("runtime error: first name cannot be nil")
	}
	if lastName == nil {
		panic("runtime error: last name cannot be nil")
	}
	fmt.Printf("%s %s\n", *firstName, *lastName)
	fmt.Println("returned normally from fullName")
}
func CheckRecover() {
	fmt.Println("Recover Process Pending")
	if r := recover(); r != nil {
		fmt.Println("Catch Panic ==>", r)
	}
	fmt.Printf("Recover Process Finish")
}
func Panic() {
	lastName := "Musk"
	defer CheckRecover()     //Recover phát hiện Panic và đưa thông báo kết thúc
	fullName(nil, &lastName) //Panic =>> Crash
	fmt.Println("returned normally from main")
}
