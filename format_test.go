package main

import "testing"

func TestECLSpecifier(t *testing.T) {
	tests := []struct {
		name string
		ecl  ErrorCorrectionLevel
		want int
	}{
		{
			name: "error correction is L",
			ecl:  Low,
			want: 0b01,
		},
		{
			name: "error correction is M",
			ecl:  Medium,
			want: 0b00,
		},
		{
			name: "error correction is Q",
			ecl:  High,
			want: 0b11,
		},
		{
			name: "error correction is H",
			ecl:  Highest,
			want: 0b10,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ecls := ECLSpecifier(test.ecl)
			if ecls != test.want {
				t.Errorf("expected %d, got %d\n", test.want, ecls)
			}
		})
	}
}
