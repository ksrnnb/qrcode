package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ksrnnb/qrcode"
)

func main() {
	content := flag.String("c", "Hello, World!", "content of qrcode")

	flag.Parse()

	s, err := qrcode.New(qrcode.ECL_Medium, *content)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot be encoded: %v\n", err)
		os.Exit(1)
	}
	p, err := s.PNG(255)
	if err != nil {
		fmt.Fprintf(os.Stderr, "png encode error: %v\n", err)
		os.Exit(1)
	}
	err = os.WriteFile("qrcode.png", p, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "write file error: %v\n", err)
		os.Exit(1)
	}
}
