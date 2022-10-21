package main

import "fmt"

func main() {
	// use only Medium to simplify
	ecl := ECLIndicator("M")
	modules := [][]bool{{}}
	formatInfo := FormatInfo(ecl, modules)
	fmt.Println(formatInfo)

	// EncodeAlphaNumericString("Hello")
}
