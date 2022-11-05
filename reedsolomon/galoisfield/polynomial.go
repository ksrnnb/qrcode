package galoisfield

import (
	"github.com/ksrnnb/qrcode/bitset"
)

// Polynomial means polynomial over GF(2^8)
type Polynomial struct {
	terms []Element
}

func NewPolynomial(bs *bitset.BitSet) Polynomial {
	totalBytes := bs.Length() / 8
	if bs.Length()%8 != 0 {
		totalBytes++
	}
	poly := Polynomial{
		terms: make([]Element, totalBytes),
	}
	for i := 0; i < totalBytes; i++ {
		poly.terms[i] = Element(bs.ByteAt(i))
	}
	return poly
}

func NewMonomial(e Element, degree int) Polynomial {
	if e.IsZero() {
		return Polynomial{}
	}
	m := Polynomial{
		terms: make([]Element, degree+1),
	}
	m.terms[degree] = e
	return m
}

// Add returns sum f(x) + g(x)
func (f Polynomial) Add(g Polynomial) Polynomial {
	sumMaxDegree := f.maxDegree()
	if sumMaxDegree < g.maxDegree() {
		sumMaxDegree = g.maxDegree()
	}

	sumPoly := Polynomial{
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

// Add returns product f(x) * g(x)
func (f Polynomial) Multiply(g Polynomial) Polynomial {
	fMaxDegree := f.maxDegree()
	gMaxDegree := g.maxDegree()

	product := Polynomial{
		terms: make([]Element, fMaxDegree+gMaxDegree+1),
	}

	for fi := 0; fi <= fMaxDegree; fi++ {
		for gi := 0; gi <= gMaxDegree; gi++ {
			if f.terms[fi] == 0 || g.terms[gi] == 0 {
				continue
			}
			p := NewMonomial(f.terms[fi].Multiply(g.terms[gi]), fi+gi)
			product = product.Add(p)
		}
	}
	return product.normalize()
}

// Remainder returns remainder of f(x) / g(x)
func (f Polynomial) Remainder(g Polynomial) Polynomial {
	if g.IsZero() {
		panic("polynomial cannot be divided by zero")
	}

	remainder := f
	for remainder.maxDegree() >= g.maxDegree() {
		// calculate quotient
		qDeg := remainder.maxDegree() - g.maxDegree()
		qCoff := remainder.maxDegreeTerm().Divide(g.maxDegreeTerm())
		q := NewMonomial(Element(qCoff), qDeg)

		remainder = remainder.Add(g.Multiply(q))
	}
	return remainder
}

// IsZero returns true if polynomial has no terms
func (f Polynomial) IsZero() bool {
	for _, term := range f.terms {
		if !term.IsZero() {
			return false
		}
	}
	return true
}

func (f Polynomial) maxDegreeTerm() Element {
	return f.terms[f.maxDegree()]
}

// maxDegree returns max degree of polynomial
func (f Polynomial) maxDegree() int {
	return len(f.terms) - 1
}

// normalize returns new polynomial which is normalized
// if term of max degree is zero, it will be removed
func (f Polynomial) normalize() Polynomial {
	maxDegree := f.maxDegree()
	newMaxDegree := maxDegree
	for i := maxDegree; i >= 0; i-- {
		if f.terms[i] != 0 {
			break
		}
		newMaxDegree--
	}
	if newMaxDegree < 0 {
		return Polynomial{}
	}
	f.terms = f.terms[0 : newMaxDegree+1]
	return f
}
