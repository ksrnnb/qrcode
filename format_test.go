package main

import (
	"testing"

	"github.com/ksrnnb/qrcode/bitset"
)

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

func TestCharacterCountIndicatorBits(t *testing.T) {
	tests := []struct {
		name    string
		version int
		mode    ModeIndicator
		want    int
	}{
		{
			name:    "version is less than minimum value",
			version: 0,
			mode:    Numeric,
			want:    0,
		},
		{
			name:    "version is greater than maximum value",
			version: 41,
			mode:    Numeric,
			want:    0,
		},
		{
			name:    "version is 1-9 and mode is numeric",
			version: 1,
			mode:    Numeric,
			want:    10,
		},
		{
			name:    "version is 1-9 and mode is alpha numeric",
			version: 1,
			mode:    AlphaNumeric,
			want:    9,
		},
		{
			name:    "version is 1-9 and mode is 8 bits byte",
			version: 1,
			mode:    EightBits,
			want:    8,
		},
		{
			name:    "version is 1-9 and mode is kanji",
			version: 1,
			mode:    Kanji,
			want:    8,
		},
		{
			name:    "version is 10-26 and mode is numeric",
			version: 10,
			mode:    Numeric,
			want:    12,
		},
		{
			name:    "version is 10-26 and mode is alpha numeric",
			version: 10,
			mode:    AlphaNumeric,
			want:    11,
		},
		{
			name:    "version is 10-26 and mode is 8 bits byte",
			version: 10,
			mode:    EightBits,
			want:    16,
		},
		{
			name:    "version is 10-26 and mode is kanji",
			version: 10,
			mode:    Kanji,
			want:    10,
		},
		{
			name:    "version is 27-40 and mode is numeric",
			version: 27,
			mode:    Numeric,
			want:    14,
		},
		{
			name:    "version is 27-40 and mode is alpha numeric",
			version: 27,
			mode:    AlphaNumeric,
			want:    13,
		},
		{
			name:    "version is 27-40 and mode is 8 bits byte",
			version: 27,
			mode:    EightBits,
			want:    16,
		},
		{
			name:    "version is 27-40 and mode is kanji",
			version: 27,
			mode:    Kanji,
			want:    12,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := characterCountIndicatorBits(test.version, test.mode)
			if result != test.want {
				t.Errorf("expected %d, got %d\n", test.want, result)
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

func TestAddZeroPadding(t *testing.T) {
	tests := []struct {
		name     string
		codeSize int
		length   int
		arg      int
		wantPos  int
	}{
		{
			name:     "equals to code length",
			codeSize: 2,
			length:   16,
			arg:      0b1111_1111_1111_1111,
			wantPos:  16,
		},
		{
			name:     "less than code length and greater than code length-4",
			codeSize: 2,
			length:   13,
			arg:      0b1_1111_1111_1111,
			wantPos:  16,
		},
		{
			name:     "less than code code length-4",
			codeSize: 2,
			length:   10,
			arg:      0b11_1111_1111,
			wantPos:  14,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bs := bitset.NewBitSet(test.codeSize * 8)
			nextPos := bs.SetInt(0, test.arg, test.length)
			nextPos = addZeroPadding(bs, nextPos, test.codeSize*8)
			if nextPos != test.wantPos {
				t.Errorf("nextPos is expected %d, but got %d\n", test.wantPos, nextPos)
			}
		})
	}
}
