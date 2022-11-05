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

// Values returns byte array of bitset for debug
func (bs *BitSet) Values() []byte {
	length := bs.length / 8
	if length%8 != 0 {
		length++
	}
	values := make([]byte, length)
	for i := 0; i < length; i++ {
		values[i] = bs.ByteAt(i)
	}
	return values
}

func (bs *BitSet) SetInt(v int, length int) {
	bs.ensureCapacity(length)
	for i := 0; i < length; i++ {
		bs.value[bs.pos+i] = GetBit(v, length-1-i)
	}
	bs.pos += length
}

func (bs *BitSet) SetBytes(bytes []byte) {
	for _, v := range bytes {
		bs.SetByte(v)
	}
}

func (bs *BitSet) SetByte(v byte) {
	bs.SetInt(int(v), 8)
}

func (bs *BitSet) SetBool(v bool) {
	bs.ensureCapacity(1)
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

func (bs *BitSet) ByteAt(bytePos int) byte {
	if bytePos < 0 || bytePos >= bs.length {
		panic("index is invalid")
	}

	bytePos *= 8
	var v byte
	for i := bytePos; i < bytePos+8 && i < bs.length; i++ {
		v <<= 1
		if bs.GetValue(i) {
			v |= 1
		}
	}
	return v
}

func (bs *BitSet) Position() int {
	return bs.pos
}

func (bs *BitSet) Length() int {
	return bs.length
}

func (bs *BitSet) Clone() *BitSet {
	newbs := &BitSet{
		length: bs.length,
		value:  make([]bool, bs.length),
		pos:    0,
	}
	newbs.SetBools(bs.value...)
	return newbs
}

func (bs *BitSet) ensureCapacity(num int) {
	if bs.pos+num <= bs.length {
		return
	}
	lack := bs.pos + num - bs.length
	bs.value = append(bs.value, make([]bool, lack)...)
	bs.length += lack
}
