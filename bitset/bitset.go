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

func (bs *BitSet) SetInt(pos int, v int, size int) {
	for i := 0; i < size; i++ {
		bs.value[pos+i] = v>>(size-1-i) == 1
	}
}

func (bs *BitSet) SetByte(pos int, v byte) {
	bs.SetInt(pos, int(v), 8)
}

func (bs *BitSet) GetValue(pos int) bool {
	return bs.value[pos]
}
