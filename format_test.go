package main

import (
	"testing"
)

func TestFormatInfo(t *testing.T) {
	tests := []struct {
		name string
		ecl  ErrorCorrectionLevel
		want uint16
	}{
		{
			name: "error correction level is M and mask pattern is 101",
			ecl:  ECL_Medium,
			want: 0b1000_0001_1001_110,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := FormatInfo(test.ecl, [][]bool{{}})
			if result != test.want {
				t.Errorf("expected %b, got %b\n", test.want, result)
			}
		})
	}
}
