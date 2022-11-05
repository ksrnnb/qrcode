package reedsolomon

import "github.com/ksrnnb/qrcode/reedsolomon/galoisfield"

func GeneratorPolynomial(degree int) galoisfield.Polynomial {
	if degree < 2 {
		panic("degree must be over 2")
	}

	generator := galoisfield.NewMonomial(galoisfield.Element(1), 0)

	for i := 0; i < degree; i++ {
		elem := galoisfield.ElementByExponentOfAlpha(i)
		root := galoisfield.NewMonomial(elem, 0)
		x := galoisfield.NewMonomial(1, 1)

		generator = generator.Multiply(x.Add(root))
	}
	return generator
}
