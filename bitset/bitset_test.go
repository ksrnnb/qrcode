package bitset

import "testing"

func TestSetByte(t *testing.T) {
	var v uint8 = 0b10111010
	bs := NewBitSet(8)

	bs.SetByte(0, v)

	for i := 0; i < 8; i++ {
		want := v>>(7-i) == 1
		result := bs.GetValue(i)
		if want != result {
			t.Errorf("expected %v, got %v at pos %d\n", want, result, i)
		}
	}
}
