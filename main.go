package main

import (
	"fmt"

	"github.com/ksrnnb/qrcode/reedsolomon"
)

func main() {
	// use only Medium to simplify
	data := EncodeRawData(ECL_Medium, "Hello, World!")
	// 1-M: count of error collection words is 10
	ecwords := 10
	rsEncoded := reedsolomon.Encode(data, ecwords)

	fmt.Printf("%+v\n", rsEncoded.Values())

	// TODO: add regular symbol, evaluate mask,
	modules := [][]bool{{}}
	formatInfo := FormatInfo(ECL_Medium, modules)
	fmt.Println(formatInfo)

	// EncodeAlphaNumericString("Hello")
}
