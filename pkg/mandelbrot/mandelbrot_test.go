package mandelbrot_test

import (
	"testing"

	mdb "github.com/bwarren2/mandelbrot/pkg/mandelbrot"
)

func TestDraw(t *testing.T) {
	mdb.Draw("test_img.png", 90, 45, -2.5, 1, -1, 1)
}

func TestScale(t *testing.T) {

	testcases := []struct {
		name                                                    string
		inputMin, inputMax, outputMin, outputMax, input, output float64
	}{
		{"Rescale 1-2 to 3-4, 1->3", 1, 2, 3, 4, 1, 3},
		{"Rescale 1-2 to 3-4, 2->4", 1, 2, 3, 4, 2, 4},
		{"Rescale 1-2 to 3-4, 1.25->3.25", 1, 2, 3, 4, 1.25, 3.25},
		{"Rescale 2-1 to 3-4, 2->3", 2, 1, 3, 4, 2, 3},
		{"Rescale 2-1 to 3-4, 1->4", 2, 1, 3, 4, 1, 4},
		{"Rescale 2-1 to 3-4, 1.5->3.5", 2, 1, 3, 4, 1.5, 3.5},
	}
	for _, tc := range testcases {
		fn := mdb.Scale(tc.inputMin, tc.inputMax, tc.outputMin, tc.outputMax)
		value := fn(tc.input)
		if value != tc.output {
			t.Fatalf("got %v for %v", value, tc.name)
		}
	}
}
