package bitset

func GetBit[T int | uint8 | uint16](v T, pos int) bool {
	return ((v >> pos) & 1) == 1
}

type BitSet struct {
	length int
	value  []bool
	pos    int
}

func NewBitSet(length int) *BitSet {
	return &BitSet{
		length: length,
		value:  make([]bool, length),
		pos:    0,
	}
}

func (bs *BitSet) SetInt(v int, length int) {
	for i := 0; i < length; i++ {
		bs.value[bs.pos+i] = GetBit(v, length-1-i)
	}
	bs.pos += length
}

func (bs *BitSet) SetByte(v byte) {
	bs.SetInt(int(v), 8)
}

func (bs *BitSet) SetBool(v bool) {
	bs.value[bs.pos] = v
	bs.pos++
}

func (bs *BitSet) SetBools(v ...bool) {
	for _, b := range v {
		bs.SetBool(b)
	}
}

func (bs *BitSet) GetValue(pos int) bool {
	return bs.value[pos]
}

func (bs *BitSet) Position() int {
	return bs.pos
}

func (bs *BitSet) Length() int {
	return bs.length
}
