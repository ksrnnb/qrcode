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

func FormatInfo(ecl ErrorCorrectionLevel, modules [][]bool) uint16 {
	mask := EvaluateMask(modules)
	formatBitSequence := (uint8(ecl) << 3) | mask

	return maskedBitSequence[formatBitSequence]
}
