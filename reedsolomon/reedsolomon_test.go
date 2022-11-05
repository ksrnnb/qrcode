package reedsolomon

import (
	"testing"

	"github.com/ksrnnb/qrcode/bitset"
	"github.com/ksrnnb/qrcode/reedsolomon/galoisfield"
)

func TestEncode(t *testing.T) {
	var tests = []struct {
		name    string
		ecwords int
		data    []byte
		want    []byte
	}{
		{
			name:    "count of error correction words is 5",
			ecwords: 5,
			data: []byte{
				0b01000000, 0b00011000, 0b10101100, 0b11000011, 0b00000000,
			},
			want: []byte{
				0b01000000, 0b00011000, 0b10101100, 0b11000011, 0b00000000,
				0b10000110, 0b00001101, 0b00100010, 0b10101110, 0b00110000,
			},
		},
		{
			name:    "count of error correction words is 10",
			ecwords: 10,
			data: []byte{
				0b00010000, 0b00100000, 0b00001100, 0b01010110, 0b01100001, 0b10000000, 0b11101100, 0b00010001, 0b11101100, 0b00010001, 0b11101100, 0b00010001, 0b11101100, 0b00010001, 0b11101100, 0b00010001,
			},
			want: []byte{
				0b00010000, 0b00100000, 0b00001100, 0b01010110, 0b01100001, 0b10000000, 0b11101100, 0b00010001, 0b11101100, 0b00010001, 0b11101100, 0b00010001, 0b11101100, 0b00010001, 0b11101100, 0b00010001,
				0b10100101, 0b00100100, 0b11010100, 0b11000001, 0b11101101, 0b00110110, 0b11000111, 0b10000111, 0b00101100, 0b01010101,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bs := bitset.NewBitSet(len(test.data) * 8)
			bs.SetBytes(test.data)
			result := Encode(bs, test.ecwords)
			for i, want := range test.want {
				if result.ByteAt(i) != want {
					t.Errorf("want %d, but got %d at index: %d\n", want, result.ByteAt(i), i)
					break
				}
			}
		})
	}
}

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
