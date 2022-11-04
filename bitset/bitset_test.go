package bitset

import "testing"

func TestGetBit(t *testing.T) {
	v := 0b1011101
	bools := []bool{true, false, true, true, true, false, true}

	for i, want := range bools {
		result := GetBit(v, i)
		if result != want {
			t.Errorf("expected %v, but got %v at pos %d\n", want, result, i)
		}
	}
}

func TestSetInt(t *testing.T) {
	v := 0b101110101111001000
	codeLength := 18
	bs := NewBitSet(codeLength)

	bs.SetInt(v, codeLength)

	for i := 0; i < codeLength; i++ {
		want := GetBit(v, codeLength-1-i)
		result := bs.GetValue(i)
		if want != result {
			t.Errorf("expected %v, got %v at pos %d\n", want, result, i)
		}
	}
}

func TestSetByte(t *testing.T) {
	var v uint8 = 0b10111010
	codeLength := 8
	bs := NewBitSet(codeLength)

	bs.SetByte(v)

	for i := 0; i < codeLength; i++ {
		want := GetBit(v, codeLength-1-i)
		result := bs.GetValue(i)
		if want != result {
			t.Errorf("expected %v, got %v at pos %d\n", want, result, i)
		}
	}
}

func TestBools(t *testing.T) {
	bools := []bool{true, false, true}
	bs := NewBitSet(len(bools))
	bs.SetBools(bools...)
	for i, want := range bools {
		result := bs.GetValue(i)
		if result != want {
			t.Errorf("expected %v, but got %v", want, result)
		}
	}
}

func TestByteAt(t *testing.T) {
	bools := []bool{
		true, true, false, false, true, true, true, false,
		false, true, false, false, true, false, false, true,
	}

	tests := []struct {
		pos  int
		want byte
	}{
		{
			pos:  0,
			want: 206,
		},
		{
			pos:  1,
			want: 73,
		},
	}

	for _, test := range tests {
		t.Run("test for ByteAt", func(t *testing.T) {
			bs := NewBitSet(len(bools))
			bs.SetBools(bools...)
			if bs.ByteAt(test.pos) != test.want {
				t.Errorf("want %d, but got %d", test.want, bs.ByteAt(test.pos))
			}
		})
	}
}
