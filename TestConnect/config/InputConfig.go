package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const StringInvalid = "Blank String! Please try again..."
const NumberInvalid = "Number format invalid or Out of Range! Please try again..."
const BoolInvalid = "Input Bool invalid! Please try again..."

func input() string {
	result := ""
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		result = scanner.Text()
	}
	return result
}

func InputString() string {
	input := input()
	result := strings.Trim(input, " ")
	if result == "" {
		fmt.Println(StringInvalid)
		return InputString()
	}
	return result
}

func InputInteger() int {
	input := input()
	value, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		fmt.Println(NumberInvalid)
		return InputInteger()
	}
	return int(value)
}

func InputFloat() float64 {
	fmt.Println(NumberInvalid)
	input := input()
	value, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return InputFloat()
	}
	return value
}

func InputBool() bool {
	input := input()
	value, err := strconv.ParseBool(input)
	if err != nil {
		fmt.Println(BoolInvalid)
		return InputBool()
	}
	return value
}
