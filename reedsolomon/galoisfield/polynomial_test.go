package galoisfield

import (
	"testing"
)

func TestPolynomial_Add(t *testing.T) {
	tests := []struct {
		name          string
		f             Polynomial
		g             Polynomial
		want          Polynomial
		wantMaxDegree int
	}{
		{
			name: "normal addition",
			f: Polynomial{
				terms: []Element{
					0b0011_0011, 0b0110_0110, 0b1100_1100,
				},
			},
			g: Polynomial{
				terms: []Element{
					0b0000_0001, 0b0000_0010, 0b0000_0100, 0b0000_1000,
				},
			},
			want: Polynomial{
				terms: []Element{
					0b0011_0010, 0b0110_0100, 0b1100_1000, 0b0000_1000,
				},
			},
		},
		{
			name: "max degree becomes zero will be normalized",
			f: Polynomial{
				terms: []Element{
					0b0011_0011, 0b0110_0110, 0b1100_1100,
				},
			},
			g: Polynomial{
				terms: []Element{
					0b0000_0001, 0b0000_0010, 0b1100_1100,
				},
			},
			want: Polynomial{
				terms: []Element{
					0b0011_0010, 0b0110_0100,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sum := test.f.Add(test.g)
			if sum.maxDegree() != test.want.maxDegree() {
				t.Errorf("want degree is %d, but got %d\n", test.want.maxDegree(), sum.maxDegree())
				return
			}

			for i, v := range test.want.terms {
				if sum.terms[i] != v {
					t.Errorf("want %+v, but got %+v\n", test.want, sum)
					break
				}
			}

		})
	}
}

func TestPolynomial_Multiply(t *testing.T) {
	tests := []struct {
		name          string
		f             Polynomial
		g             Polynomial
		want          Polynomial
		wantMaxDegree int
	}{
		{
			name: "normal multiply",
			f: Polynomial{
				terms: []Element{
					0b1000_0000, 0b1000_0111, // α^7 + α^13*x^1
				},
			},
			g: Polynomial{
				terms: []Element{
					0b0111_0100, 0b0001_0000, 0b0001_1101, // α^10 + α^4*x^1 + α^8*x^2
				},
			},
			want: Polynomial{
				terms: []Element{
					0b1001_1000, 0b0010_0001, 0b1011_1110, 0b0111_0101, // α^17 + α^138*x^1 + α^65*x^2 + α^21*x^3
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := test.f.Multiply(test.g)
			if p.maxDegree() != test.want.maxDegree() {
				t.Errorf("want degree is %d, but got %d\n", test.want.maxDegree(), p.maxDegree())
				return
			}

			for i, v := range test.want.terms {
				if p.terms[i] != v {
					t.Errorf("want %+v, but got %+v\n", test.want, p)
					break
				}
			}

		})
	}
}

func TestPolynomial_Remainder(t *testing.T) {
	tests := []struct {
		name          string
		f             Polynomial
		g             Polynomial
		want          Polynomial
		wantMaxDegree int
	}{
		{
			name: "normal remainder",
			f: Polynomial{
				terms: []Element{
					0b0111_0100, 0b0001_0000, 0b0001_1101, // α^10 + α^4*x^1 + α^8*x^2
				},
			},
			g: Polynomial{
				terms: []Element{
					0b1000_0000, 0b1000_0111, // α^7 + α^13*x^1
				},
			},
			want: Polynomial{
				terms: []Element{
					0b1110_1011, // α^235
				},
			},
		},
		{
			name: "remainder is zero",
			f: Polynomial{
				terms: []Element{
					0b0111_0100, 0b0001_0000, 0b0001_1101, // α^10 + α^4*x^1 + α^8*x^2
				},
			},
			g: Polynomial{
				terms: []Element{
					0b0111_0100, 0b0001_0000, 0b0001_1101, // α^10 + α^4*x^1 + α^8*x^2
				},
			},
			want: zeroPolynomial,
		},
		{
			name: "2",
			f:    Polynomial{terms: []Element{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1}},
			g:    Polynomial{terms: []Element{1, 1, 1, 0, 1, 1, 0, 0, 1, 0, 1}},
			want: Polynomial{terms: []Element{0, 0, 1, 1, 1, 0, 1, 1}},
		},
		{
			name: "3",
			f:    Polynomial{terms: []Element{91, 50, 25, 184, 194, 105, 45, 244, 58, 44}},
			g:    Polynomial{terms: []Element{254, 120, 88, 44, 11, 1}},
			want: Polynomial{terms: []Element{}},
		},
		{
			name: "4",
			f:    Polynomial{terms: []Element{0, 0, 0, 0, 0, 0, 195, 172, 24, 64}},
			g:    Polynomial{terms: []Element{116, 147, 63, 198, 31, 1}},
			want: Polynomial{terms: []Element{48, 174, 34, 13, 134}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := test.f.Remainder(test.g)
			if p.maxDegree() != test.want.maxDegree() {
				t.Errorf("want degree is %d, but got %d\n", test.want.maxDegree(), p.maxDegree())
				return
			}

			for i, v := range test.want.terms {
				if p.terms[i] != v {
					t.Errorf("want %+v, but got %+v\n", test.want, p)
					break
				}
			}
		})
	}
}
