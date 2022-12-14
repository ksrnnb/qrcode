package qrcode

import (
	"unicode/utf8"

	"github.com/ksrnnb/qrcode/bitset"
	"github.com/ksrnnb/qrcode/reedsolomon"
)

func encodeRawData(info qrInfo) (*bitset.BitSet, error) {
	// 8 bytes mode => countDataCodeWords * 8
	codeLength := info.countDataCodeWords * 8
	bs := bitset.NewBitSet(codeLength)

	bitCount := characterCountIndicatorBits(info.version, info.mode)
	charCount := utf8.RuneCountInString(info.src)

	addModeIndicator(bs, info.mode)
	addCharacterCountIndicator(bs, bitCount, charCount)
	addSrcData(bs, info.src)
	addTerminator(bs)
	addPaddingBit(bs)

	rsEncoded := reedsolomon.Encode(bs, info.countErrorCordWords())

	return rsEncoded, nil
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
