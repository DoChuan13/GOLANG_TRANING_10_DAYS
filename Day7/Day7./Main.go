package main

import (
	"Day7/Practice/ReadWrite"
	"fmt"
)

func main() {
	fmt.Println("======================1. JSON Example======================")
	//json.TestJsonExample()
	fmt.Println("======================2.0. Read File======================")
	//ReadWrite.ReadFile(ReadWrite.FilePath)
	fmt.Println("======================2.1. Write File (OverWrite)======================")
	//Overwrite new content to files
	//ReadWrite.OverWriteFile(ReadWrite.FilePath)
	fmt.Println("======================2.2. Write File (Continues)======================")
	file := ReadWrite.OpenFile(ReadWrite.FilePath)
	input := []string{"Hello", "World"}
	ReadWrite.WriteFile(file, input)
	fmt.Println("======================3. Testing======================")
	fmt.Println("======================4. Command Line======================")
	//CommandLine.Command()
}
