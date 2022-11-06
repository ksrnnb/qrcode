package main

import (
	"fmt"
	"math"

	"github.com/ksrnnb/qrcode/reedsolomon"
)

func main() {
	// use only Medium to simplify
	ecl := ECL_Medium
	data := EncodeRawData(ecl, "Hello, World!")

	// 1-M: count of error collection words is 10
	ecwords := 10
	rsEncoded := reedsolomon.Encode(data, ecwords)

	var s *Symbol
	penalty := math.MaxInt
	for mask := uint8(0b000); mask <= uint8(0b111); mask++ {
		newS := NewSymbol(ecl, mask, rsEncoded)
		newS.build()
		if newS.penalty() < penalty {
			penalty = newS.penalty()
			s = newS
		}
	}
	fmt.Printf("%+v\n", s)
}
