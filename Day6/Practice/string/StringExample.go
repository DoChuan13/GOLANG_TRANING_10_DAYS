package string

import (
	"fmt"
	"strings"
)

func StringExample() {
	var str = "   Hello World   "
	fmt.Println("Original String ==>", str)
	fmt.Println("Trim Function", strings.Trim(str, " "))
	fmt.Println("Prefix==>", strings.HasPrefix(str, " "))
	fmt.Println("Prefix==>", strings.HasSuffix(str, " "))
	fmt.Println("Contain==>", strings.Contains(str, "Hello"))
	fmt.Println("Contain Any==>", strings.ContainsAny(str, "Hlo"))
}
