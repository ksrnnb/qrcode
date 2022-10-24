package main

import (
	"unicode/utf8"

	"github.com/ksrnnb/qrcode/bitset"
)

type ErrorCorrectionLevel uint8

const (
	// Level L
	Low ErrorCorrectionLevel = 0b01

	// Level M
	Medium ErrorCorrectionLevel = 0b00

	// Level Q
	High ErrorCorrectionLevel = 0b11

	// Level H
	Highest ErrorCorrectionLevel = 0b10
)

const (
	modeCharCount = 4
)

// maskedBitSequence means masking (5, 15, 7) BCH code
// reference: JIS X0510 : 2018 (ISO/IEC 18004 : 2015) Table C.1
var maskedBitSequence = []uint16{
	0x5412,
	0x5125,
	0x5E7C,
	0x5B4B,
	0x45F9,
	0x40CE,
	0x4F97,
	0x4AA0,
	0x77C4,
	0x72F3,
	0x7DAA,
	0x789D,
	0x662F,
	0x6318,
	0x6C41,
	0x6976,
	0x1689,
	0x13BE,
	0x1CE7,
	0x19D0,
	0x0762,
	0x0255,
	0x0D0C,
	0x083B,
	0x355F,
	0x3068,
	0x3F31,
	0x3A06,
	0x24B4,
	0x2183,
	0x2EDA,
	0x2BED,
}

// ECLIndicator returns error correction level indicator
func ECLIndicator(level string) ErrorCorrectionLevel {
	switch level {
	case "L":
		return Low
	case "M":
		return Medium
	case "Q":
		return High
	case "H":
		return Highest
	default:
		return Medium
	}
}

type ModeIndicator uint8

// reference: JIS X0510 : 2018 (ISO/IEC 18004 : 2015) Table 2
const (
	Numeric ModeIndicator = 1 << iota

	AlphaNumeric

	// 8 bits byte
	EightBits

	Kanji
)

func FormatInfo(ecl ErrorCorrectionLevel, modules [][]bool) uint16 {
	mask := EvaluateMask(modules)
	formatBitSequence := (uint8(ecl) << 3) | mask

	return maskedBitSequence[formatBitSequence]
}

func EncodeRawData(ecl ErrorCorrectionLevel, src string) *bitset.BitSet {
	if ecl != Medium {
		panic("this app supports only 1-M type, error correction level must be 'M'")
	}

	if utf8.RuneCountInString(src) >= 16 {
		panic("this app supports only 1-M type and 8 bits byte mode, must be less than 16 characters")
	}

	mode := EightBits
	version := 1

	// 1-M: data code size is 16
	// 8 bytes mode => 16 * 8
	codeLength := 16 * 8
	bs := bitset.NewBitSet(codeLength)

	bitCount := characterCountIndicatorBits(version, mode)
	charCount := utf8.RuneCountInString(src)

	addModeIndicator(bs, mode)
	addCharacterCountIndicator(bs, bitCount, charCount)
	addSrcData(bs, src)
	// TODO: bitset should have current position and Length() method
	addTerminator(bs)
	addPaddingBit(bs)

	return bs
}

// characterCountIndicatorBits returns character count indicater's bit numbers
// reference: JIS X0510 : 2018 (ISO/IEC 18004 : 2015) Table 3
func characterCountIndicatorBits(version int, mode ModeIndicator) int {
	if 1 <= version && version <= 9 {
		return version1To9CharacterCountIndicatorBits(mode)
	} else if 10 <= version && version <= 26 {
		return version10To26CharacterCountIndicatorBits(mode)
	} else if 27 <= version && version <= 40 {
		return version27To40CharacterCountIndicatorBits(mode)
	}
	return 0
}

func version1To9CharacterCountIndicatorBits(mode ModeIndicator) int {
	switch mode {
	case Numeric:
		return 10
	case AlphaNumeric:
		return 9
	case EightBits:
		return 8
	case Kanji:
		return 8
	default:
		return 0
	}
}

func version10To26CharacterCountIndicatorBits(mode ModeIndicator) int {
	switch mode {
	case Numeric:
		return 12
	case AlphaNumeric:
		return 11
	case EightBits:
		return 16
	case Kanji:
		return 10
	default:
		return 0
	}
}

func version27To40CharacterCountIndicatorBits(mode ModeIndicator) int {
	switch mode {
	case Numeric:
		return 14
	case AlphaNumeric:
		return 13
	case EightBits:
		return 16
	case Kanji:
		return 12
	default:
		return 0
	}
}

// addModeIndicator adds mode indicator and returns next position
func addModeIndicator(bs *bitset.BitSet, mode ModeIndicator) {
	bs.SetInt(int(mode), modeCharCount)
}

// addCharacterCountIndicator adds character count indicator and returns next position
func addCharacterCountIndicator(bs *bitset.BitSet, bitCount int, charCount int) {
	bs.SetInt(charCount, bitCount)
}

// addSrcData adds src data and returns next position
func addSrcData(bs *bitset.BitSet, src string) {
	// supports only 8 bit byte mode
	for _, c := range src {
		bs.SetByte(byte(c))
	}
}

// addTerminator adds 0000 padding and returns next position
func addTerminator(bs *bitset.BitSet) {
	if bs.Position() == bs.Length() {
		return
	}

	if bs.Position() <= bs.Length()-4 {
		bs.SetInt(0, 4)
		return
	}

	nextPos := bs.Position()
	for i := nextPos; i < bs.Length(); i++ {
		bs.SetBool(false)
	}
}

// addZeroPadding add 0 padding if last bit string is not 8 bits
func addZeroPadding(bs *bitset.BitSet) {
	for i := 0; i < 8; i++ {
		if bs.Position()%8 == 0 {
			break
		}
		bs.SetBool(false)
	}
}

// addPaddingBit add 0 padding if last bit string is not 8 bits
func addPaddingBit(bs *bitset.BitSet) {
	addZeroPadding(bs)
	if bs.Position() == bs.Length() {
		return
	}

	paddingPatterns := []int{0b11101100, 00010001}
	for i := 0; bs.Position() < bs.Length(); i++ {
		bs.SetInt(paddingPatterns[i%2], 8)
	}
}
