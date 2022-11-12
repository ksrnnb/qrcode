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

func TestSetBytes(t *testing.T) {
	tests := []struct {
		name   string
		values []byte
	}{
		{
			name:   "appending bytes",
			values: []byte{0b0101_0101, 0b1111_1111, 0b0000_1111},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bs := NewBitSet(0)
			bs.SetBytes(test.values)
			for i, want := range test.values {
				if bs.ByteAt(i) != want {
					t.Errorf("want %b, but got %b at index: %d", want, bs.ByteAt(i), i)
					break
				}
			}
		})
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
	tests := []struct {
		pos  int
		data []bool
		want byte
	}{
		{
			pos: 0,
			data: []bool{
				true, true, false, false, true, true, true, false,
				false, true, false, false, true, false, false, true,
			},
			want: 206,
		},
		{
			pos: 1,
			data: []bool{
				true, true, false, false, true, true, true, false,
				false, true, false, false, true, false, false, true,
			},
			want: 73,
		},
		{
			pos: 0,
			data: []bool{
				false, true, false, false, true, false,
			},
			want: 18,
		},
	}

	for _, test := range tests {
		t.Run("test for ByteAt", func(t *testing.T) {
			bs := NewBitSet(len(test.data))
			bs.SetBools(test.data...)
			if bs.ByteAt(test.pos) != test.want {
				t.Errorf("want %d, but got %d", test.want, bs.ByteAt(test.pos))
			}
		})
	}
}

func TestClone(t *testing.T) {
	tests := []struct {
		name   string
		values []bool
	}{
		{
			name:   "Clone",
			values: []bool{false, true, false, false, true, true, true},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bs := NewBitSet(len(test.values))
			bs.SetBools(test.values...)
			c := bs.Clone()
			for i, want := range test.values {
				if c.value[i] != want {
					t.Errorf("want %v, but got %v at index: %d\n", want, c.value[i], i)
					break
				}
			}
		})
	}
}
