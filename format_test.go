package main

import "testing"

func TestECLIndicator(t *testing.T) {
	tests := []struct {
		name string
		ecl  string
		want ErrorCorrectionLevel
	}{
		{
			name: "error correction is L",
			ecl:  "L",
			want: Low,
		},
		{
			name: "error correction is M",
			ecl:  "M",
			want: Medium,
		},
		{
			name: "error correction is Q",
			ecl:  "Q",
			want: High,
		},
		{
			name: "error correction is H",
			ecl:  "H",
			want: Highest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ecls := ECLIndicator(test.ecl)
			if ecls != test.want {
				t.Errorf("expected %d, got %d\n", test.want, ecls)
			}
		})
	}
}

func TestFormatInfo(t *testing.T) {
	tests := []struct {
		name string
		ecl  ErrorCorrectionLevel
		want uint16
	}{
		{
			name: "error correction level is M and mask pattern is 101",
			ecl:  Medium,
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
