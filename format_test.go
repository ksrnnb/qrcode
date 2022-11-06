package main

import (
	"testing"
)

func TestFormatInfo(t *testing.T) {
	tests := []struct {
		name string
		ecl  ErrorCorrectionLevel
		mask uint8
		want []bool
	}{
		{
			name: "error correction level is M and mask pattern is 101",
			ecl:  ECL_Medium,
			mask: 0b101,
			want: []bool{
				true, false, false, false,
				false, false, false, true,
				true, false, false, true,
				true, true, false,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := FormatInfo(test.ecl, test.mask)
			for i, want := range test.want {
				if result.GetValue(i) != want {
					t.Errorf("expected %v, got %v at index: %d\n", want, result.GetValue(i), i)
				}
			}
		})
	}
}
