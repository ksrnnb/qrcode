package bitset

type BitSet struct {
	size  int
	value []bool
}

func NewBitSet(size int) *BitSet {
	return &BitSet{
		size:  size,
		value: make([]bool, size),
	}
}

func GetBit[T int | uint8 | uint16](v T, pos int) bool {
	return ((v >> pos) & 1) == 1
}

func (bs *BitSet) SetInt(pos int, v int, size int) (nextPos int) {
	for i := 0; i < size; i++ {
		bs.value[pos+i] = GetBit(v, size-1-i)
	}
	return pos + size
}

func (bs *BitSet) SetByte(pos int, v byte) (nextPos int) {
	return bs.SetInt(pos, int(v), 8)
}

func (bs *BitSet) SetBool(pos int, v bool) (nextPos int) {
	bs.value[pos] = v
	return pos + 1
}

func (bs *BitSet) SetBools(pos int, v ...bool) (nextPos int) {
	for i, b := range v {
		nextPos = bs.SetBool(pos+i, b)
	}
	return nextPos
}

func (bs *BitSet) GetValue(pos int) bool {
	return bs.value[pos]
}
