package galoisfield

import (
	"testing"

	"github.com/ksrnnb/qrcode/bitset"
)

func TestAdd_Polynominal(t *testing.T) {
	tests := []struct {
		name          string
		f             []bool
		g             []bool
		want          []bool
		wantMaxDegree int
	}{
		{
			name: "normal addition",
			f: []bool{
				false, false, true, true, false, false, true, true,
				false, true, true, false, false, true, true, false,
				true, true, false, false, true, true, false, false,
			},
			g: []bool{
				false, false, false, false, false, false, false, true,
				false, false, false, false, false, false, true, false,
				false, false, false, false, false, true, false, false,
				false, false, false, false, true, false, false, false,
			},
			want: []bool{
				false, false, true, true, false, false, true, false,
				false, true, true, false, false, true, false, false,
				true, true, false, false, true, false, false, false,
				false, false, false, false, true, false, false, false,
			},
			wantMaxDegree: 3,
		},
		{
			name: "max degree becomes zero will be normalized",
			f: []bool{
				false, false, true, true, false, false, true, true,
				false, true, true, false, false, true, true, false,
				true, true, false, false, true, true, false, false,
			},
			g: []bool{
				false, false, false, false, false, false, false, true,
				false, false, false, false, false, false, true, false,
				true, true, false, false, true, true, false, false,
			},
			want: []bool{
				false, false, true, true, false, false, true, false,
				false, true, true, false, false, true, false, false,
			},
			wantMaxDegree: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := bitset.NewBitSet(len(test.f))
			g := bitset.NewBitSet(len(test.g))
			want := bitset.NewBitSet(len(test.want))
			f.SetBools(test.f...)
			g.SetBools(test.g...)
			want.SetBools(test.want...)
			fPoly := NewPolynominal(f)
			gPoly := NewPolynominal(g)
			wantPoly := NewPolynominal(want)
			sum := fPoly.Add(gPoly)

			for i, v := range wantPoly.terms {
				if sum.terms[i] != v {
					t.Errorf("want %+v, but got %+v\n", wantPoly, sum)
					break
				}
			}
			if sum.maxDegree() != test.wantMaxDegree {
				t.Errorf("want degree is %d, but go %d\n", test.wantMaxDegree, sum.maxDegree())
			}
		})
	}
}
