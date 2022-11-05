package reedsolomon

import (
	"testing"

	"github.com/ksrnnb/qrcode/reedsolomon/galoisfield"
)

func TestGeneratorPolynomial(t *testing.T) {
	tests := []struct {
		name     string
		degree   int
		elements []galoisfield.Element
	}{
		{
			name:   "normal remainder",
			degree: 10,
			elements: []galoisfield.Element{
				// α^45+α^32x +α^94x^2 +α^64x^3 +α^70x^4 +α^118x^5 +α^61x^6 +α^46x^7 +α^67x^8 +α^251x^9 + x^10
				0b1100_0001, 0b1001_1101, 0b0111_0001, 0b0101_1111, 0b0101_1110, 0b1100_0111, 0b0110_1111, 0b1001_1111, 0b1100_0010, 0b1101_1000, 0b0000_0001,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			generator := GeneratorPolynomial(test.degree)
			gTerms := generator.Terms()
			for i, want := range test.elements {
				if gTerms[i] != want {
					t.Errorf("want %d, but got %d\n", want, gTerms[i])
				}
			}
		})
	}
}
