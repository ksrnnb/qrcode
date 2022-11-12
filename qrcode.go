package qrcode

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"unicode/utf8"

	"github.com/ksrnnb/qrcode/bitset"
)

type QRCode struct {
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

	separatorHorizontalPattern = [][]bool{
		{false, false, false, false, false, false, false, false},
	}
	separatorVerticalPattern = [][]bool{
		{false},
		{false},
		{false},
		{false},
		{false},
		{false},
		{false},
		{false},
	}
)

type qrInfo struct {
	version            int
	ecl                ErrorCorrectionLevel
	mode               ModeIndicator
	dataCap            int // code cap = countDataCodeWords + countErrorCodeWords
	countDataCodeWords int
	srcCap             int
	src                string
}

func newQRInfo(ecl ErrorCorrectionLevel, src string) qrInfo {
	// supports only version 1
	switch ecl {
	case ECL_Low:
		return qrInfo{
			version:            1,
			ecl:                ecl,
			mode:               EightBits,
			dataCap:            26,
			countDataCodeWords: 19,
			srcCap:             17,
			src:                src,
		}
	case ECL_Medium:
		return qrInfo{
			version:            1,
			ecl:                ecl,
			mode:               EightBits,
			dataCap:            26,
			countDataCodeWords: 16,
			srcCap:             14,
			src:                src,
		}
	case ECL_High:
		return qrInfo{
			version:            1,
			ecl:                ecl,
			mode:               EightBits,
			dataCap:            26,
			countDataCodeWords: 13,
			srcCap:             11,
			src:                src,
		}
	default: // Error Correction Level: H
		return qrInfo{
			version:            1,
			ecl:                ecl,
			mode:               EightBits,
			dataCap:            26,
			countDataCodeWords: 9,
			srcCap:             7,
			src:                src,
		}
	}
}

func (qi qrInfo) countErrorCordWords() int {
	return qi.dataCap - qi.countDataCodeWords
}

func New(ecl ErrorCorrectionLevel, content string) (*QRCode, error) {
	info := newQRInfo(ecl, content)
	if utf8.RuneCountInString(content) > info.srcCap {
		return nil, fmt.Errorf("this app supports only version 1 and 8 bits byte mode, must be less than %d characters", info.srcCap+1)
	}

	data, err := encodeRawData(info)
	if err != nil {
		return nil, err
	}

	var q *QRCode
	penalty := math.MaxInt
	for mask := uint8(0b000); mask <= uint8(0b111); mask++ {
		newQR := newQRCode(ecl, mask, data)
		if newQR.penalty() < penalty {
			penalty = newQR.penalty()
			q = newQR
		}
	}
	return q, nil
}

func newQRCode(ecl ErrorCorrectionLevel, mask uint8, data *bitset.BitSet) *QRCode {
	// version1: module size per line is 21
	size := 21

	q := &QRCode{
		ecl:     ecl,
		mask:    mask,
		data:    data,
		modules: make([][]bool, size+2*quietZoneSize),
		dirties: make([][]bool, size+2*quietZoneSize),
		size:    size,
	}

	for i := range q.modules {
		q.modules[i] = make([]bool, size+2*quietZoneSize)
		q.dirties[i] = make([]bool, size+2*quietZoneSize)
	}

	q.build()

	return q
}

func (q *QRCode) Image(size int) image.Image {
	realSize := q.size + 2*quietZoneSize

	if size < realSize {
		size = realSize
	}

	// Output image.
	rect := image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{size, size}}

	// Saves a few bytes to have them in this order
	p := color.Palette([]color.Color{color.White, color.Black})
	img := image.NewPaletted(rect, p)
	fgClr := uint8(img.Palette.Index(color.Black))

	// QR code bitmap.
	bitmap := q.modules

	// Map each image pixel to the nearest QR code module.
	modulesPerPixel := float64(realSize) / float64(size)
	for y := 0; y < size; y++ {
		y2 := int(float64(y) * modulesPerPixel)
		for x := 0; x < size; x++ {
			x2 := int(float64(x) * modulesPerPixel)

			v := bitmap[y2][x2]

			if v {
				pos := img.PixOffset(x, y)
				img.Pix[pos] = fgClr
			}
		}
	}

	return img
}

func (q *QRCode) PNG(size int) ([]byte, error) {
	img := q.Image(size)

	var b bytes.Buffer
	err := png.Encode(&b, img)

	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (q *QRCode) build() {
	q.addFinderPatterns()
	q.addSeparatorPattern()
	q.addTimingPatterns()
	// NOTE: format info is added after applying mask on JIS 7.1 section
	//       but format info is added before applying mask here because dirties should be marked before adding data
	q.addFormatInfo()
	q.addData()
}

func (q *QRCode) addFinderPatterns() {
	// top left
	q.add2dPattern(0, 0, finderPattern)

	// top right
	q.add2dPattern(q.size-finderPatternSize, 0, finderPattern)

	// bottom left
	q.add2dPattern(0, q.size-finderPatternSize, finderPattern)
}

func (q *QRCode) addSeparatorPattern() {
	// top left vertical
	q.add2dPattern(finderPatternSize, 0, separatorVerticalPattern)
	// top left horizontal
	q.add2dPattern(0, finderPatternSize, separatorHorizontalPattern)

	// top right vertical
	q.add2dPattern(q.size-finderPatternSize-1, 0, separatorVerticalPattern)
	// top right horizontal
	q.add2dPattern(q.size-finderPatternSize-1, finderPatternSize, separatorHorizontalPattern)

	// bottom left vertical
	q.add2dPattern(finderPatternSize, q.size-finderPatternSize-1, separatorVerticalPattern)
	// bottom left horizontal
	q.add2dPattern(0, q.size-finderPatternSize-1, separatorHorizontalPattern)
}

func (q *QRCode) addTimingPatterns() {
	// timing pattern starts with true
	v := true

	// start of timing pattern: finder pattern size + separator size (1)
	for i := finderPatternSize + 1; i < q.size-finderPatternSize-1; i++ {
		// horizontal direction
		q.add(i, finderPatternSize-1, v)
		// vertical direction
		q.add(finderPatternSize-1, i, v)
		// next module is inverse boolean
		v = !v
	}
}

func (q *QRCode) addData() {
	// when dx is  0, position is right
	// when dx is -1, position is left
	dx := 0

	// start from bottom right
	x := q.size - 1
	y := q.size - 1

	// direction
	direction := up

	for i := 0; i < q.data.Length(); i++ {
		mask := calculateMask(x+dx, y, q.mask)
		// != is equivalent to XOR.
		q.add(x+dx, y, mask != q.data.GetValue(i))

		if i == q.data.Length()-1 {
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
					if y < q.size-1 {
						y++
					} else {
						// if y is bottom, change direction
						direction = up
						x -= 2
					}
				}
			}

			// column 6 cannot be write and need to skip
			if x == 6 {
				x--
			}

			if !q.isDirty(x+dx, y) {
				// break if next position is not dirty
				break
			}
			// if next position is dirty, tries to find next not dirty position
		}
	}
}

func (q *QRCode) addFormatInfo() {
	fi := FormatInfo(q.ecl, q.mask)
	q.addVerticalFormatInfo(fi)
	q.addHorizontalFormatInfo(fi)
}

func (q *QRCode) addVerticalFormatInfo(fi *bitset.BitSet) {
	last := formatInfoLength - 1
	// Bits 0-5
	for i := 0; i <= 5; i++ {
		q.add(finderPatternSize+1, i, fi.GetValue(last-i))
	}

	// (x, y) = (finderPatternSize+1, 6) is ignored, because it is timing pattern

	// Bits 6-7
	for i := 6; i <= 7; i++ {
		q.add(finderPatternSize+1, i+1, fi.GetValue(last-i))
	}

	// (finderPatternSize+1, q.size-finderPatternSize-1) is black
	q.add(finderPatternSize+1, q.size-finderPatternSize-1, true)

	// Bits 8-14
	for i := 8; i <= 14; i++ {
		q.add(finderPatternSize+1, q.size-finderPatternSize-8+i, fi.GetValue(last-i))
	}
}

func (q *QRCode) addHorizontalFormatInfo(fi *bitset.BitSet) {
	last := formatInfoLength - 1
	// Bits 0-7
	for i := 0; i <= 7; i++ {
		q.add(q.size-i-1, finderPatternSize+1, fi.GetValue(last-i))
	}

	// Bits 8
	q.add(finderPatternSize, finderPatternSize+1, fi.GetValue(last-8))

	// (x, y) = (finderPatternSize-1, finderPatternSize+1) is ignored, because it is timing pattern

	// Bits 9-14
	for i := 9; i <= 14; i++ {
		q.add(14-i, finderPatternSize+1, fi.GetValue(last-i))
	}
}

func (q *QRCode) penalty() int {
	return q.penalty1() + q.penalty2() + q.penalty3() + q.penalty4()
}

func (q *QRCode) penalty1() int {
	p := q.penalty1Horizontal()
	if p < q.penalty1Vertical() {
		p = q.penalty1Vertical()
	}

	return p
}

func (q *QRCode) penalty1Horizontal() int {
	penalty := 0
	penaltyWeight := 3

	for y := 0; y < q.size; y++ {
		lastValue := q.get(0, y)
		count := 1

		for x := 1; x < q.size; x++ {
			v := q.get(x, y)

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

func (q *QRCode) penalty1Vertical() int {
	penalty := 0
	penaltyWeight := 3

	for x := 0; x < q.size; x++ {
		lastValue := q.get(x, 0)
		count := 1

		for y := 1; y < q.size; y++ {
			v := q.get(x, y)

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

func (q *QRCode) penalty2() int {
	penalty := 0
	penaltyWeight2 := 3

	for y := 1; y < q.size; y++ {
		for x := 1; x < q.size; x++ {
			topLeft := q.get(x-1, y-1)
			above := q.get(x, y-1)
			left := q.get(x-1, y)
			current := q.get(x, y)

			if current == left && current == above && current == topLeft {
				penalty++
			}
		}
	}
	return penalty * penaltyWeight2
}

func (q *QRCode) penalty3() int {
	penaltyWeight3 := 40

	for y := 0; y < q.size; y++ {
		var bitBuffer uint16 = 0x00

		for x := 0; x < q.size; x++ {
			bitBuffer <<= 1
			if v := q.get(x, y); v {
				bitBuffer |= 1
			}

			switch bitBuffer & 0x7ff {
			// 0b000 0101 1101 or 0b101 1101 0000
			// 0x05d           or 0x5d0
			case 0x05d, 0x5d0:
				return penaltyWeight3
			default:
				if x == q.size-1 && (bitBuffer&0x7f) == 0x5d {
					return penaltyWeight3
				}
			}
		}
	}

	for x := 0; x < q.size; x++ {
		var bitBuffer uint16 = 0x00

		for y := 0; y < q.size; y++ {
			bitBuffer <<= 1
			if v := q.get(x, y); v {
				bitBuffer |= 1
			}

			switch bitBuffer & 0x7ff {
			// 0b000 0101 1101 or 0b101 1101 0000
			// 0x05d           or 0x5d0
			case 0x05d, 0x5d0:
				return penaltyWeight3
			default:
				if y == q.size-1 && (bitBuffer&0x7f) == 0x5d {
					return penaltyWeight3
				}
			}
		}
	}

	return 0
}

func (q *QRCode) penalty4() int {
	penaltyWeight4 := 10
	numModules := q.size * q.size
	numDarkModules := 0

	for x := 0; x < q.size; x++ {
		for y := 0; y < q.size; y++ {
			if v := q.get(x, y); v {
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

func (q *QRCode) add2dPattern(x int, y int, pattern [][]bool) {
	for dy, row := range pattern {
		for dx, v := range row {
			q.add(x+dx, y+dy, v)
		}
	}
}

func (q *QRCode) add(x int, y int, v bool) {
	q.modules[y+quietZoneSize][x+quietZoneSize] = v
	q.dirties[y+quietZoneSize][x+quietZoneSize] = true
}

func (q *QRCode) get(x int, y int) bool {
	return q.modules[y+quietZoneSize][x+quietZoneSize]
}

func (q *QRCode) isDirty(x, y int) bool {
	return q.dirties[y+quietZoneSize][x+quietZoneSize]
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
