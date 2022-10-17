package main

type ErrorCorrectionLevel int

const (
	// Level L
	Low ErrorCorrectionLevel = iota + 1

	// Level M
	Medium

	// Level Q
	High

	// Level H
	Highest
)

// ECLSpecifier returns error correction level specifier
func ECLSpecifier(ecl ErrorCorrectionLevel) int {
	switch ecl {
	case Low:
		return 0b01
	case Medium:
		return 0b00
	case High:
		return 0b11
	case Highest:
		return 0b10
	default:
		return 0b00
	}
}
