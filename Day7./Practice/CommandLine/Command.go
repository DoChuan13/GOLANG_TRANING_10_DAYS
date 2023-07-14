package CommandLine

import (
	"fmt"
	"os"
)

func Command() {
	argsWithProg := os.Args
	argsWithoutProg := os.Args[1:]
	fre := os.Args[0]
	arg := os.Args[3]
	fmt.Println("Full Main exe ==>", fre)
	fmt.Println("Pre Arg==>", argsWithProg)
	fmt.Println("Arg==>", argsWithoutProg)
	fmt.Println("Arg[3]==>", arg)
}
