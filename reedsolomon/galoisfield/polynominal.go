package galoisfield

import (
	"github.com/ksrnnb/qrcode/bitset"
)

// Polynominal means polynominal over GF(2^8)
type Polynominal struct {
	terms []Element
}

func NewPolynominal(bs *bitset.BitSet) Polynominal {
	totalBytes := bs.Length() / 8
	if bs.Length()%8 != 0 {
		totalBytes++
	}
	poly := Polynominal{
		terms: make([]Element, totalBytes),
	}
	for i := 0; i < totalBytes; i++ {
		poly.terms[i] = Element(bs.ByteAt(i))
	}
	return poly
}

// Add returns sum of polynominal over GF(2^8)
func (f Polynominal) Add(g Polynominal) Polynominal {
	sumMaxDegree := f.maxDegree()
	if sumMaxDegree < g.maxDegree() {
		sumMaxDegree = g.maxDegree()
	}

	sumPoly := Polynominal{
		terms: make([]Element, sumMaxDegree+1),
	}

	for i := 0; i <= sumMaxDegree; i++ {
		if i <= f.maxDegree() && i <= g.maxDegree() {
			sumPoly.terms[i] = f.terms[i].Add(g.terms[i])
		} else if i <= f.maxDegree() {
			sumPoly.terms[i] = f.terms[i]
		} else {
			sumPoly.terms[i] = g.terms[i]
		}
	}
	return sumPoly.normalize()
}

// maxDegree returns max degree of polynominal
func (f Polynominal) maxDegree() int {
	return len(f.terms) - 1
}

// normalize returns new polynominal which is normalized
// if term of max degree is zero, it will be removed
func (f Polynominal) normalize() Polynominal {
	maxDegree := f.maxDegree()
	newMaxDegree := maxDegree
	for i := maxDegree; i >= 0; i-- {
		if f.terms[i] != 0 {
			break
		}
		newMaxDegree--
	}
	if newMaxDegree < 0 {
		return Polynominal{}
	}
	f.terms = f.terms[0 : newMaxDegree+1]
	return f
}
