# qrcode

This pacakge can simple QR code encoding. It only supoprts version 1, so very short string can only be encoded.

# How to use

## install package

```bash
go get github.com/ksrnnb/qrcode
```

## Example

```go
package main

import (
	"fmt"
	"os"

	"github.com/ksrnnb/qrcode"
)

func main() {
	q, err := qrcode.New(qrcode.ECL_Medium, "Hello, World")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot be encoded: %v\n", err)
		return
	}

	size := 255
	p, err := q.PNG(size)
	if err != nil {
		fmt.Fprintf(os.Stderr, "png encode error: %v\n", err)
		return
	}

	err = os.WriteFile("qrcode.png", p, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "write file error: %v\n", err)
		return
	}
}
```

# Reference

- https://github.com/skip2/go-qrcode

# Japanese Referenct
- JIS X 0510 2004
  - https://kikakurui.com/x0/X0510-2004-01.html (without image)
  - https://www.jisc.go.jp/app/jis/general/GnrJISNumberNameSearchList?toGnrJISStandardDetailList (only view)
- [ＱＲコードをつくってみる](https://www.swetake.com/qrcode/qr1.html)
- [例題で学ぶ符号理論入門](https://www.morikita.co.jp/books/mid/081741)
- [例題が語る符号理論 BCH符号・RS符号・QRコード](https://www.kyoritsu-pub.co.jp/book/b10010558.html)
