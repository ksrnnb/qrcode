package main

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

func FormatInfo(ecl ErrorCorrectionLevel, modules [][]bool) uint16 {
	mask := EvaluateMask(modules)
	formatBitSequence := (uint8(ecl) << 3) | mask

	return maskedBitSequence[formatBitSequence]
}
