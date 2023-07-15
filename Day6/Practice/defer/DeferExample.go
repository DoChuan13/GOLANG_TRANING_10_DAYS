package _defer

import "fmt"

type Person struct {
	name string
}

func getNameDefer(name *Person) {
	fmt.Println("Person info (Df)  ==>", name)
}
func Defer() {
	a := 5
	fmt.Println("Original value of a===>", a)
	defer fmt.Println("Value defer of a===>", a) //Value a = 5, câu lệnh thực hiện sau cùng
	a = 10
	fmt.Println("Value of a after changed===>", a)
	//Person
	user := Person{name: "Chuan"}
	fmt.Println("User Info (Or)===>", user)
	defer getNameDefer(&user) //Name = "Vuong", câu lệnh thực hiện cận cuối (Stack) giá trị bị thay đổi do truyền con trỏ
	user.name = "Vuong"
	fmt.Println("User Info (After)===>", user)
}
