package main

import (
	"Demo/other"
	"fmt"
	"reflect"
)

func main() {
	fmt.Println("==========================Variable=============================")
	//*. Variable Declaration
	var checkValue = 3 > 10 //false
	println("Result===>", checkValue)
	var pi0 float64
	fmt.Println(pi0)
	var pi1 float64 = 3.14
	fmt.Println(pi1)
	var pi2 = 3.14 // same pi2:=3.14
	fmt.Println(pi2)
	var name1, name2 string = "Chuan", "Vuong"
	fmt.Println("Multiple Variable Declaration===>", name1, name2)

	//*Ép kiểu dữ liệu
	i := 55
	j := 67.8
	//sum := i + j  //(lỗi)

	sum1 := i + int(j) //(OK)
	fmt.Println("Ép kiểu dữ liệu==>", sum1)

	fmt.Println("==========================Access Modifier=============================")
	//Access Modifier
	//Other Package
	fmt.Println(other.Other1)
	other.OtherFunc()

	//Same Package
	fmt.Println(same0)
	fmt.Println(Same1)
	sameFunc()
	SameFunc()

	var x int = 12
	var y int32 = 12
	//*Constant (Untyped Constant/ Typed Constant)
	fmt.Println("==========================Const=============================")
	const strConst = 12
	fmt.Println(strConst == x)
	fmt.Println(strConst == y)
	fmt.Println("Number =>> ", strConst, " with type ", reflect.TypeOf(strConst))
	const string0 = "Hello World,おはようございます"
	fmt.Println("String =>> ", string0, " with type ", reflect.TypeOf(string0))
	const string1 string = "Hello world"
	fmt.Println("String =>> ", string1, " with type ", reflect.TypeOf(string1))

	//For Loop
	fmt.Println("==========================For Loop=============================")
	for i := 1; i <= 10; i++ {
		fmt.Printf(" %d", i)
	}

	//If Else
	fmt.Println("\n==========================If Statement=============================")
	var ifNum = 10
	if ifNum%2 == 0 {
		fmt.Println("Even")
	}
	fmt.Println("\n==========================Switch Case 1=============================")
	switchValue := 0
	switch switchValue {
	case 0:
		fmt.Println("In Zero Value")
	case 1, 2:
		fmt.Println("Multiple matches")
	case 42: // Don't "fall through".
		fmt.Println("reached")
	case 43:
		fmt.Println("Unreached")
	default:
		fmt.Println("Optional")
	}

	fmt.Println("\n==========================Switch Case 2=============================")
	month := 5
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		fmt.Printf("Month %d has 31 days", month)
	case 4, 6, 9, 11:
		fmt.Printf("Month %d has 30 days", month)
	case 2:
		fmt.Printf("Month %d has 28 or 29 days", month)
	default:
		fmt.Printf("Invalid")
	}
	fmt.Println("\n==========================Switch Case 3=============================")
	switch month {
	case 1:
		fallthrough
	case 3:
		fallthrough
	case 5:
		fallthrough
	case 7:
		fallthrough
	case 8:
		fallthrough
	case 10:
		fallthrough
	case 12:
		fmt.Printf("Month %d has 31 days", month)
	case 4:
		fallthrough
	case 6:
		fallthrough
	case 9:
		fallthrough
	case 11:
		fmt.Printf("Month %d has 30 days", month)
	case 2:
		fmt.Printf("Month %d has 28 or 29 days", month)
	default:
		fmt.Printf("Invalid")
	}
}
