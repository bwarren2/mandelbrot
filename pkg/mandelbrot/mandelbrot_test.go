package mandelbrot_test

import (
	"encoding/gob"
	"image"
	"os"
	"testing"

	"github.com/go-test/deep"

	"github.com/bwarren2/mandelbrot/pkg/mandelbrot"
	mdb "github.com/bwarren2/mandelbrot/pkg/mandelbrot"
)

// TestGif _wants_ to test creating a small mandelbrot gif, but can;t encode a sample
func TestGif(t *testing.T) {
	colors := mandelbrot.NewPalette(10)
	mdb.Gif(10, 10, 3, 10, -1.5, 0, .98, colors)
	// var want image.RGBA
	// f, _ := os.Open("sample_gif.dat") // Is there a cleaner way to do this/
	// defer f.Close()
	// gob.Register(gif.GIF)
	// enc := gob.NewEncoder(f)
	// err := enc.Encode(got)
	// if err != nil {
	// 	log.Fatal("encode error:", err)
	// }
}

// TestDraw tests creating a small mandelbrot image and compares to a known-good image
func TestDraw(t *testing.T) {
	colors := mandelbrot.NewPalette(10)
	got := mdb.Draw(10, 5, 10, -2.5, 1, -1, 1, colors)
	var want image.RGBA
	f, _ := os.Open("../../test/sample.dat") // Is there a cleaner way to do this/
	defer f.Close()
	dec := gob.NewDecoder(f)
	err := dec.Decode(&want)
	if err != nil {
		t.Errorf("decode error: %v", err)
	}
	if diff := deep.Equal(got, &want); diff != nil { // Is there a better way to do this without deep?
		t.Error(diff)
	}
}

// TestScale tests the scale-generation function, demonstrating it scales and clamps output
func TestScale(t *testing.T) {
	testcases := []struct {
		name                                                    string
		inputMin, inputMax, outputMin, outputMax, input, output float64
	}{
		{"Rescale 1-2 to 3-4, 1->3", 1, 2, 3, 4, 1, 3},
		{"Rescale 1-2 to 3-4, 2->4", 1, 2, 3, 4, 2, 4},
		{"Rescale 1-2 to 3-4, 1.25->3.25", 1, 2, 3, 4, 1.25, 3.25},
		{"Rescale 1-2 to 3-4, 1->3", 1, 2, 3, 4, 0, 3},
		{"Rescale 1-2 to 3-4, 2->4", 1, 2, 3, 4, 5, 4},
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
