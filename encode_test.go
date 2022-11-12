package qrcode

import (
	"testing"

	"github.com/ksrnnb/qrcode/bitset"
)

func TestEncodeRawData(t *testing.T) {
	tests := []struct {
		name string
		ecl  ErrorCorrectionLevel
		data string
		want []byte
	}{
		{
			name: "1-M encode",
			ecl:  ECL_Medium,
			data: "Hello, World!",
			want: []byte{64, 212, 134, 86, 198, 198, 242, 194, 5, 118, 247, 38, 198, 66, 16, 236},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			info := newQRInfo(test.ecl, test.data)
			result, err := encodeRawData(info)
			if err != nil {
				t.Errorf("error: %v\n", err)
				return
			}
			for i, want := range test.want {
				if result.ByteAt(i) != want {
					t.Errorf("want %d, but got %d at index: %d\n", want, result.ByteAt(i), i)
					break
				}
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

func TestAddTerminator(t *testing.T) {
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
			bs.SetInt(test.arg, test.length)
			addTerminator(bs)
			if bs.Position() != test.wantPos {
				t.Errorf("nextPos is expected %d, but got %d\n", test.wantPos, bs.Position())
			}
		})
	}
}

func TestAddZeroPadding(t *testing.T) {
	tests := []struct {
		name   string
		bsSize int
		length int
		arg    int
		want   int
	}{
		{
			name:   "last bit string is 8bits",
			bsSize: 16,
			length: 16,
			arg:    0b1111_1111_1111_1111,
			want:   0b1111_1111_1111_1111,
		},
		{
			name:   "last bit string is not 8bits",
			bsSize: 16,
			length: 13,
			arg:    0b1_1111_1111_1111,
			want:   0b1111_1111_1111_1000,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bs := bitset.NewBitSet(test.bsSize)
			bs.SetInt(test.arg, test.length)
			addZeroPadding(bs)

			for i := 0; i < test.bsSize; i++ {
				result := bs.GetValue(i)
				want := (test.want>>(test.bsSize-1-i))&1 == 1
				if result != want {
					t.Errorf("expected %v, but got %v at pos %d\n", want, result, i)
				}
			}
		})
	}
}
