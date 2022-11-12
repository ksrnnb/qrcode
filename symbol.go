package main

import (
	"math"

	"github.com/ksrnnb/qrcode/bitset"
)

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

	up   = 1
	down = 2
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

	for i := range s.modules {
		s.modules[i] = make([]bool, size+2*quietZoneSize)
		s.dirties[i] = make([]bool, size+2*quietZoneSize)
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
	// when dx is  0, position is right
	// when dx is -1, position is left
	dx := 0

	// start from bottom right
	x := s.size - 1
	y := s.size - 1

	// direction
	direction := up

	for i := 0; i < s.data.Length(); i++ {
		mask := calculateMask(x+dx, y, s.mask)
		// != is equivalent to XOR.
		s.add(x+dx, y, mask != s.data.GetValue(i))

		if i == s.data.Length()-1 {
			break
		}

		for {
			if dx == 0 {
				// next position is left
				dx = -1
			} else {
				// next position is right
				dx = 0

				if direction == up {
					if y > 0 {
						y--
					} else {
						// if y is top, change direction
						direction = down
						x -= 2
					}
				} else {
					if y < s.size-1 {
						y++
					} else {
						// if y is bottom, change direction
						direction = up
						x -= 2
					}
				}
			}

			if !s.isDirty(x+dx, y) {
				// break if next position is not dirty
				break
			}
			// if next position is dirty, tries to find next not dirty position
		}
	}
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
	p := s.penalty1Horizontal()
	if p < s.penalty1Vertical() {
		p = s.penalty1Vertical()
	}

	return p
}

func (s *Symbol) penalty1Horizontal() int {
	penalty := 0
	penaltyWeight := 3

	for y := 0; y < s.size; y++ {
		lastValue := s.get(0, y)
		count := 1

		for x := 1; x < s.size; x++ {
			v := s.get(x, y)

			if v != lastValue {
				count = 1
				lastValue = v
			} else {
				count++
				if count == 5 {
					if penalty < penaltyWeight {
						penalty = penaltyWeight
					}
				} else if count > 6 {
					if penalty < penaltyWeight+count-5 {
						penalty = penaltyWeight + count - 5
					}
				}
			}
		}
	}
	return penalty
}

func (s *Symbol) penalty1Vertical() int {
	penalty := 0
	penaltyWeight := 3

	for x := 0; x < s.size; x++ {
		lastValue := s.get(x, 0)
		count := 1

		for y := 1; y < s.size; y++ {
			v := s.get(x, y)

			if v != lastValue {
				count = 1
				lastValue = v
			} else {
				count++
				if count == 5 {
					if penalty < penaltyWeight {
						penalty = penaltyWeight
					}
				} else if count > 6 {
					if penalty < penaltyWeight+count-5 {
						penalty = penaltyWeight + count - 5
					}
				}
			}
		}
	}
	return penalty
}

func (s *Symbol) penalty2() int {
	penalty := 0
	penaltyWeight2 := 3

	for y := 1; y < s.size; y++ {
		for x := 1; x < s.size; x++ {
			topLeft := s.get(x-1, y-1)
			above := s.get(x, y-1)
			left := s.get(x-1, y)
			current := s.get(x, y)

			if current == left && current == above && current == topLeft {
				penalty++
			}
		}
	}
	return penalty * penaltyWeight2
}

func (s *Symbol) penalty3() int {
	penaltyWeight3 := 40

	for y := 0; y < s.size; y++ {
		var bitBuffer uint16 = 0x00

		for x := 0; x < s.size; x++ {
			bitBuffer <<= 1
			if v := s.get(x, y); v {
				bitBuffer |= 1
			}

			switch bitBuffer & 0x7ff {
			// 0b000 0101 1101 or 0b101 1101 0000
			// 0x05d           or 0x5d0
			case 0x05d, 0x5d0:
				return penaltyWeight3
			default:
				if x == s.size-1 && (bitBuffer&0x7f) == 0x5d {
					return penaltyWeight3
				}
			}
		}
	}

	for x := 0; x < s.size; x++ {
		var bitBuffer uint16 = 0x00

		for y := 0; y < s.size; y++ {
			bitBuffer <<= 1
			if v := s.get(x, y); v {
				bitBuffer |= 1
			}

			switch bitBuffer & 0x7ff {
			// 0b000 0101 1101 or 0b101 1101 0000
			// 0x05d           or 0x5d0
			case 0x05d, 0x5d0:
				return penaltyWeight3
			default:
				if y == s.size-1 && (bitBuffer&0x7f) == 0x5d {
					return penaltyWeight3
				}
			}
		}
	}

	return 0
}

func (s *Symbol) penalty4() int {
	penaltyWeight4 := 10
	numModules := s.size * s.size
	numDarkModules := 0

	for x := 0; x < s.size; x++ {
		for y := 0; y < s.size; y++ {
			if v := s.get(x, y); v {
				numDarkModules++
			}
		}
	}

	ratio := float64(numDarkModules) / float64(numModules)
	diffPercent := 50 - ratio*100
	if diffPercent < 0 {
		diffPercent *= -1
	}

	return penaltyWeight4 * (int(math.Ceil(diffPercent / 5)))
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

func (s *Symbol) get(x int, y int) bool {
	return s.modules[y+quietZoneSize][x+quietZoneSize]
}

func (s *Symbol) isDirty(x, y int) bool {
	return s.dirties[y+quietZoneSize][x+quietZoneSize]
}

func calculateMask(x, y int, mask uint8) bool {
	// i is row position, y
	// j is column position, x
	// substitute i, j for easy comparison with JIS
	i := y
	j := x

	switch mask {
	case 0:
		return (i+j)%2 == 0
	case 1:
		return i%2 == 0
	case 2:
		return j%3 == 0
	case 3:
		return (i+j)%3 == 0
	case 4:
		return (i/2+j/3)%2 == 0
	case 5:
		return (i*j)%2+(i*j)%3 == 0
	case 6:
		return ((i*j)%2+((i*j)%3))%2 == 0
	case 7:
		return ((i+j)%2+((i*j)%3))%2 == 0
	default:
		return false
	}
}
