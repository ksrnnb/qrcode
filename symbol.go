package main

import "github.com/ksrnnb/qrcode/bitset"

type Symbol struct {
	ecl     ErrorCorrectionLevel
	mask    uint8
	data    *bitset.BitSet
	modules [][]bool
	dirties [][]bool
	size    int
}

const (
	quietZoneSize     = 4
	finderPatternSize = 7
)

var (
	finderPattern = [][]bool{
		{true, true, true, true, true, true, true},
		{true, false, false, false, false, false, true},
		{true, false, true, true, true, false, true},
		{true, false, true, true, true, false, true},
		{true, false, true, true, true, false, true},
		{true, false, false, false, false, false, true},
		{true, true, true, true, true, true, true},
	}
)

func NewSymbol(ecl ErrorCorrectionLevel, mask uint8, data *bitset.BitSet) *Symbol {
	// version1: module size per line is 21
	size := 21

	s := &Symbol{
		ecl:     ecl,
		mask:    mask,
		data:    data,
		modules: make([][]bool, size+2*quietZoneSize),
		dirties: make([][]bool, size+2*quietZoneSize),
		size:    size,
	}

	return s
}

func (s *Symbol) build() {
	s.addFinderPatterns()
	s.addTimingPatterns()
	// NOTE: format info is added after applying mask on JIS 7.1 section
	//       but format info is added before applying mask here because dirties should be marked before adding data
	s.addFormatInfo()
	s.addData()
}

func (s *Symbol) addFinderPatterns() {
	// top left
	s.add2dPattern(0, 0, finderPattern)

	// top right
	s.add2dPattern(s.size-finderPatternSize, 0, finderPattern)

	// bottom left
	s.add2dPattern(0, s.size-finderPatternSize, finderPattern)
}

func (s *Symbol) addTimingPatterns() {
	// timing pattern starts with true
	v := true

	// start of timing pattern: finder pattern size + separator size (1)
	for i := finderPatternSize + 1; i < s.size-finderPatternSize; i++ {
		// horizontal direction
		s.add(i, finderPatternSize-1, v)
		// vertical direction
		s.add(finderPatternSize-1, i, v)
		// next module is inverse boolean
		v = !v
	}
}

func (s *Symbol) addData() {
	// TODO: implement
}

func (s *Symbol) addFormatInfo() {
	fi := FormatInfo(s.ecl, s.mask)
	s.addVerticalFormatInfo(fi)
	s.addHorizontalFormatInfo(fi)
}

func (s *Symbol) addVerticalFormatInfo(fi *bitset.BitSet) {
	last := formatInfoLength - 1
	// Bits 0-5
	for i := 0; i <= 5; i++ {
		s.add(finderPatternSize+1, i, fi.GetValue(last-i))
	}

	// (x, y) = (finderPatternSize+1, 6) is ignored, because it is timing pattern

	// Bits 6-7
	for i := 6; i <= 7; i++ {
		s.add(finderPatternSize+1, i, fi.GetValue(last-i))
	}

	// (finderPatternSize+1, s.size-finderPatternSize-1) is black
	s.add(finderPatternSize+1, s.size-finderPatternSize-1, true)

	// Bits 8-14
	for i := 8; i <= 14; i++ {
		s.add(finderPatternSize+1, s.size-finderPatternSize-8+i, fi.GetValue(last-i))
	}
}

func (s *Symbol) addHorizontalFormatInfo(fi *bitset.BitSet) {
	last := formatInfoLength - 1
	// Bits 0-7
	for i := 0; i <= 7; i++ {
		s.add(s.size-i, i, fi.GetValue(last-i))
	}

	// Bits 8
	s.add(finderPatternSize, finderPatternSize+1, fi.GetValue(last-8))

	// (x, y) = (finderPatternSize-1, finderPatternSize+1) is ignored, because it is timing pattern

	// Bits 9-14
	for i := 9; i <= 14; i++ {
		s.add(14-i, finderPatternSize+1, fi.GetValue(last-i))
	}
}

func (s *Symbol) penalty() int {
	return s.penalty1() + s.penalty2() + s.penalty3() + s.penalty4()
}

func (s *Symbol) penalty1() int {
	return 0
}

func (s *Symbol) penalty2() int {
	return 0
}

func (s *Symbol) penalty3() int {
	return 0
}

func (s *Symbol) penalty4() int {
	return 0
}

func (s *Symbol) add2dPattern(x int, y int, pattern [][]bool) {
	for dy, row := range pattern {
		for dx, v := range row {
			s.add(x+dx, y+dy, v)
		}
	}
}

func (s *Symbol) add(x int, y int, v bool) {
	s.modules[y+quietZoneSize][x+quietZoneSize] = v
	s.dirties[y+quietZoneSize][x+quietZoneSize] = true
}

func (s *Symbol) isDirty(x, y int) bool {
	return s.dirties[y+quietZoneSize][x+quietZoneSize]
}
