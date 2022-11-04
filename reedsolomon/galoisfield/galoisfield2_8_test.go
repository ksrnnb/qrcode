package galoisfield

import "testing"

func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		a    Element
		b    Element
		want Element
	}{
		{
			name: "normal addition",
			a:    Element(0b0000_0010), // α^1
			b:    Element(0b0000_0100), // α^2
			want: Element(0b0000_0110), // α^26
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.a.Add(test.b) != test.want {
				t.Errorf("want %b, but got %b\n", test.want, test.a.Add(test.b))
			}
		})
	}
}

func TestSub(t *testing.T) {
	tests := []struct {
		name string
		a    Element
		b    Element
		want Element
	}{
		{
			name: "normal subtraction",
			a:    Element(0b0000_0010), // α^1
			b:    Element(0b0000_0100), // α^2
			want: Element(0b0000_0110), // α^26
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.a.Sub(test.b) != test.want {
				t.Errorf("want %b, but got %b\n", test.want, test.a.Sub(test.b))
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		name string
		a    Element
		b    Element
		want Element
	}{
		{
			name: "normal multiplication",
			a:    Element(0b0000_1000), // α^3
			b:    Element(0b0001_0000), // α^4
			want: Element(0b1000_0000), // α^7
		},
		{
			name: "multiply by 0",
			a:    Element(0b0000_1000), // α^3
			b:    Element(0b0000_0000), // 0
			want: Element(0b0000_0000), // 0
		},
		{
			name: "multiply by α^254",
			a:    Element(0b0000_1000), // α^3
			b:    Element(0b1000_1110), // α^254
			want: Element(0b0000_0100), // α^2
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.a.Multiply(test.b) != test.want {
				t.Errorf("want %b, but got %b\n", test.want, test.a.Multiply(test.b))
			}
		})
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name string
		a    Element
		b    Element
		want Element
	}{

		{
			name: "normal division",
			a:    Element(0b0001_0000), // α^4
			b:    Element(0b0000_1000), // α^3
			want: Element(0b0000_0010), // α^1
		},
		{
			name: "when index of result is minus",
			a:    Element(0b0000_1000), // α^3
			b:    Element(0b0001_0000), // α^4
			want: Element(0b1000_1110), // α^254
		},
		{
			name: "when 0 divide any element",
			a:    Element(0b0000_0000), // 0
			b:    Element(0b0000_1000), // α^3
			want: Element(0b0000_0000), // 0
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.a.Divide(test.b) != test.want {
				t.Errorf("want %b, but got %b\n", test.want, test.a.Divide(test.b))
			}
		})
	}
}
