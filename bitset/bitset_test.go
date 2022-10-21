package bitset

import "testing"

func TestSetInt(t *testing.T) {
	v := 0b101110101111001000
	size := 18
	bs := NewBitSet(size)

	bs.SetInt(0, v, size)

	for i := 0; i < size; i++ {
		want := (v >> (size - 1 - i) & 1) == 1
		result := bs.GetValue(i)
		if want != result {
			t.Errorf("expected %v, got %v at pos %d\n", want, result, i)
		}
	}
}

func TestSetByte(t *testing.T) {
	var v uint8 = 0b10111010
	size := 8
	bs := NewBitSet(size)

	bs.SetByte(0, v)

	for i := 0; i < size; i++ {
		want := (v >> (size - 1 - i) & 1) == 1
		result := bs.GetValue(i)
		if want != result {
			t.Errorf("expected %v, got %v at pos %d\n", want, result, i)
		}
	}
}

func TestBools(t *testing.T) {
	bools := []bool{true, false, true}
	bs := NewBitSet(3)
	bs.SetBools(0, bools...)
	for i, want := range bools {
		result := bs.GetValue(i)
		if result != want {
			t.Errorf("expected %v, but got %v", want, result)
		}
	}
}
