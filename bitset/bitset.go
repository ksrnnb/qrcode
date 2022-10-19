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

func (bs *BitSet) SetByte(pos int, v byte) {
	for i := 0; i < 8; i++ {
		bs.value[pos+i] = v>>(7-i) == 1
	}
}

func (bs *BitSet) GetValue(pos int) bool {
	return bs.value[pos]
}
