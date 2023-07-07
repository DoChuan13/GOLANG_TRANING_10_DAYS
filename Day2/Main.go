package main

import (
	"fmt"
	"reflect"
)

func main() {
	//1. Array
	fmt.Println("=================Array=================")
	//1. Declaration
	var numArray0 [2]int
	numArray0[0] = 10
	numArray0[1] = 20
	fmt.Println("Number Array 0=>>>> ", numArray0)

	numArray1 := [3]int{1, 2, 3}
	numArray2 := [3]int{1, 2, 3}
	numArray3 := &numArray1

	fmt.Println("Number Array 1=>>>> ", numArray1)
	fmt.Println("Number Array 2=>>>> ", numArray2)
	//So sánh giá trị hai array (tính tham trị)
	fmt.Println("Compare value of 2 arrays 1 - 2===>", numArray2 == numArray1)
	fmt.Println("Compare pointer of 2 arrays 1 - 2===>", &numArray2 == &numArray1)
	fmt.Println("Compare pointer of 2 arrays 1 - 3===>", *numArray3 == numArray1)
	numArray1[1] = 10
	fmt.Println("Number Array 1, after change=>>>> ", numArray1)
	student := [...]string{"Vuong", "Chuan", "Son"}
	fmt.Println("Student Array 0=>>>> ", student)

	numArray4 := [3]int{2: 2}
	fmt.Println(numArray4)

	//Multiple Array
	twoDirArray := [5][4]float64{
		{1, 3},
		{4.5, -3, 7.4, 2},
		{6, 2, 11},
	}

	fmt.Println("Multiple Array===>", twoDirArray)

	//
	hello := "Hello World"
	fmt.Println("String ", hello, " with length ", len(hello))
	fmt.Println(hello[0], reflect.TypeOf(hello[0]))

	//2. Slice
	fmt.Println("=================Slice=================")
	//2.1 Ref to exist array
	initArray := [5]int{1, 2, 3, 4, 5}
	fmt.Println("Initial Array ===>", initArray)
	var numberSlice []int = initArray[1:3]
	fmt.Printf("Slice value==> %d with leng %d and capacity %d\n", numberSlice, len(numberSlice), cap(numberSlice))
	initArray[2] = 100
	fmt.Println("Slice Before change Array Value", numberSlice)
	numberSlice[0] = 200
	fmt.Println("Slice Before change Slice Value", numberSlice)
	fmt.Println("Array Before change Slice Value", initArray)
	numberSlice = append(numberSlice, 300)
	fmt.Println("Slice Before append Slice Value", numberSlice)
	fmt.Println("Array Before append Slice Value", initArray)

	//2.2 create array and return slice
	cars := []string{"Ferrari", "Honda", "Ford"}
	fmt.Println("cars:", cars, "has old length", len(cars), "and capacity", cap(cars)) //capacity of cars is 3
	cars = append(cars, "Toyota")
	fmt.Println("cars:", cars, "has new length", len(cars), "and capacity", cap(cars)) //capacity of cars is doubled to 6

	//2.3 create slice with 0 value
	var nilSlice []string //zero value of a slice is nil
	fmt.Printf("Nil Slice with length %d and capacity %d", len(nilSlice), cap(nilSlice))
	fmt.Println("=================Slice - Modifier=================")
	sliceNumber := []int{1, 2, 3, 4, 5}
	fmt.Println("Original Number Slice==> ", sliceNumber)

	//Thêm phần tử vào cuối slice
	sliceNumber = append(sliceNumber, 6, 7, 8)
	fmt.Println("Append Number Slice (Last)==> ", sliceNumber)

	//Chèn phần tử vào đầu slice
	sliceNumber = append([]int{0}, sliceNumber...)
	fmt.Println("Append Number Slice (First)==> ", sliceNumber)

	//Chèn phần tử vào giữa slice
	sliceNumber = append(sliceNumber[:3], append([]int{31, 32}, sliceNumber[3:]...)...)
	fmt.Println("Append Number Slice (Mid)==> ", sliceNumber)

	//Chèn phần tử vào giữa = copy
	a := copy(sliceNumber[8:], sliceNumber[1:])
	fmt.Println(a)
	fmt.Println("Append Number Slice (Copy)==> ", sliceNumber)

	//Xóa phần tử khỏi slice
	defArr := [8]int{1, 2, 3, 4, 5, 6, 7, 8}
	deleteSlice := defArr[2:6]
	fmt.Println("Original Number Array==> ", defArr)
	fmt.Println("Original Number Slice==> ", deleteSlice, "length and cap is: ", len(deleteSlice), cap(defArr))
	deleteSlice = deleteSlice[3:]
	fmt.Println("Delete First Element Slice===>", deleteSlice)
	fmt.Println("Original Number Array==> ", defArr)

	fmt.Println("=================Map=================")
	employeeSalary := make(map[string]int)
	employeeSalary["steve"] = 12000
	employeeSalary["jamie"] = 15000
	employeeSalary["mike"] = 9000
	fmt.Println("EmployeeSalary map contents (00):", employeeSalary)

	employeeSalary1 := map[string]int{
		"steve": 12000,
		"jamie": 15000,
		"mike":  9000}
	employeeSalary1["mike"] = 12000 //modified value
	employeeSalary1["bush"] = 12000 //add new key value
	fmt.Println("EmployeeSalary map contents (01):", employeeSalary1)
	fmt.Println("Test Not Exist value in Map==>", employeeSalary1["demo"])
	delete(employeeSalary1, "mike")
	fmt.Println("EmployeeSalary map contents (01) after delete mike:", employeeSalary1)

	value, ok := employeeSalary1["mike"]
	if ok == true {
		fmt.Println("Found value is==>", value)
	} else {
		fmt.Println("Not found")
	}

	//Map with nil value
	var employeeSalary2 map[string]int
	employeeSalary2 = make(map[string]int)
	employeeSalary2["steve"] = 12000

	fmt.Println("Map with nil value after initialized===>:", employeeSalary2)
	fmt.Println("=================Slice - Modifier=================")
	forArray := []int{1, 2, 3, 4, 5, 6, 7, 8}
	for index, value := range forArray {
		fmt.Println("Index and value in For Range", index, value)
	}
}
