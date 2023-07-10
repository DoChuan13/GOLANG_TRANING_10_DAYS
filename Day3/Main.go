package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("==========================1, 2Function==========================")
	fmt.Println("Result of Function 0 =>>>> ", Add0(1, 2))
	fmt.Println("Result of Function 1 (Anonymous Original Add) =>>>> ", Add1(1, 2))

	Add1 = func(a, b int) int {
		return a * b
	}
	fmt.Println("Result of Function 1 (Anonymous Modifier to Mul) =>>>> ", Add1(2, 3))

	a, b := Swap(5, 10)
	fmt.Println("Result of Function 2 (Mul Return) =>>>> ", a, b)
	fmt.Println("Result of Function 2 (Mul Return) =>>>> ", Sum(10, []int{12, 13}...))

	fmt.Println("Before Defer Start")
	//defer DeferContext()
	//defer fmt.Println("Before Defer Context")
	fmt.Println("After Defer End")

	fmt.Println("==========================3. Closure Function==========================")
	pos, neg := adder(), adder()
	for i := 0; i < 10; i++ {
		fmt.Println(
			pos(i),
			neg(-2*i),
		)
	}
	fmt.Println("==========================4. Pointer==========================")
	var num01 = 255
	var name *int = &num01
	fmt.Println("Value of number 01 is ==>", num01)
	fmt.Println("Pointer address of number 01 (Value of pointer) is ==>", name)

	a1 := 25
	var b1 *int
	if b1 == nil {
		fmt.Println("b is", b1)
		fmt.Println("Address of b=>>>", &b)
		b1 = &a1
		fmt.Println("Pointer Address Value that b refer after initialization is", b1)
	}

	b2 := 255
	a2 := &b2
	fmt.Println("address of b is", a2)
	fmt.Println("value of b is", *a2)
	*a2++
	fmt.Println("new value of b is", b2)

	fmt.Println("==========================5. Struct==========================")
	type Student struct {
		name    string
		address string
		age     int
	}
	var student1 Student
	student1 = Student{name: "Vuong", address: "Ha Noi", age: 30}
	fmt.Println("Value of Student 01 (Struct value) =>>", student1)
	fmt.Println("Pointer Address of Student 01 (Struct value) =>>", &student1)
	student2 := Student{
		name:    "Chuan",
		address: "Nam Dinh",
		age:     31,
	}
	fmt.Println("Value of Student 02 (Struct value) =>>", student2)
	fmt.Println("Pointer Address of Student 02 (Struct value) =>>", &student2)
	var student3 Student
	fmt.Println(student3)

	type Address struct {
		city  string
		state string
	}

	type Person struct {
		name    string
		age     int
		address Address
	}
	p := Person{
		name: "Naveen",
		age:  50,
		address: Address{
			city:  "Chicago",
			state: "Illinois",
		},
	}

	fmt.Println("Name:", p.name)
	fmt.Println("Age:", p.age)
	fmt.Println("City:", p.address.city)
	fmt.Println("State:", p.address.state)

	type Person1 struct {
		name string
		age  int
		Address
	}
	p1 := Person1{
		name: "Naveen",
		age:  50,
		Address: Address{
			city:  "Chicago",
			state: "Illinois",
		},
	}

	fmt.Println("Name:", p1.name)
	fmt.Println("Age:", p1.age)
	fmt.Println("City (Promoted field):", p1.city)   //city is promoted field
	fmt.Println("State (Promoted field):", p1.state) //state is promoted field

	fmt.Println("==========================6. Method==========================")
	country := Country{id: 1, name: "USA"}
	state := State{id: 1, name: "New York"}
	country.countryInfo()
	state.countryInfo()
	fmt.Println(state)

	fmt.Println("==========================7. Interface==========================")
	var duck Animal = Duck{}
	fmt.Println("Duck scream=>>>>", duck.Scream())
	fmt.Println("Duck run with speed=>>>>", duck.Speed())

	findType("Hello")

	fmt.Println("==========================8. Error==========================")
	//Error Handing Cast value of Interface
	var ep interface{}
	ep = "Hello"
	t, ok := ep.(string)
	fmt.Println("Error Handing Hold Value of Empty Interface")
	fmt.Println(t)
	fmt.Println(ok)

	//Error Handing Cast Open Files
	file, err := os.Open("test.txt")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(file.Name(), "opened successfully")
	}

}

// All Function for Training
// hàm được đặt tên
func Add0(a, b int) int {
	return a + b
}

// hàm ẩn danh
var Add1 = func(a, b int) int {
	return a + b
}

// Nhiều tham số và nhiều giá trị trả về
func Swap(a, b int) (first, second int) { //or (int, int)
	return b, a
}

// Biến số lượng tham số 'more'
func Sum(a int, more ...int) int {
	for _, v := range more {
		a += v
	}
	return a
}

// Defer trong Function
func DeferContext() {
	fmt.Println("Content in Defer Function")
}

// Defer trong Function
func adder() func(int) int {
	sum := 0
	return func(x int) int {
		sum += x
		return sum
	}
}

type Country struct {
	id   int
	name string
}
type State struct {
	id   int
	name string
}

// Method
func (c Country) countryInfo() {
	fmt.Println("Country info =>>>", c.id, c.name)
}

func (s State) countryInfo() string {
	s.name = "Ha Noi"
	fmt.Println("State info =>>>", s.id, s.name)
	return s.name
}

// Animal Interface
type Animal interface {
	Scream() string
	Speed() float64
}

type Duck struct {
}

func (d Duck) Scream() string {
	return "Quack Quack!!!"
}

func (d Duck) Speed() float64 {
	return 5.0
}

// Empty Interface
func findType(i interface{}) {
	switch i.(type) {
	case string:
		fmt.Printf("I am a string and my value is %s\n", i.(string))
	case int:
		fmt.Printf("I am an int and my value is %d\n", i.(int))
	default:
		fmt.Printf("Unknown type\n")
	}
}
