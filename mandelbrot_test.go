package mandelbrot_test

import (
	"encoding/gob"
	"flag"
	"image"
	"image/color"
	"image/gif"
	"log"
	"os"
	"testing"

	"github.com/bwarren2/mandelbrot"
	"github.com/google/go-cmp/cmp"
)

var writeFiles = flag.Bool("write-file", false, "Write golden file")

func TestMain(m *testing.M) {
	flag.Parse()
	if *writeFiles {
		log.Print("Writing golden files")
		colors := mandelbrot.NewPalette(10)
		got := mandelbrot.MandelbrotBuilder{10, 10, 10}.Gif(3, -1.5, 0, .98, colors)
		f, err := os.Create("testdata/sample_gif.dat") // Is there a cleaner way to do this/
		if err != nil {
			panic(err)
		}
		defer f.Close()
		gob.Register(&gif.GIF{})
		gob.Register(color.RGBA{})
		enc := gob.NewEncoder(f)
		err = enc.Encode(got)
		if err != nil {
			panic(err)
		}

		img := mandelbrot.MandelbrotBuilder{10, 5, 10}.Draw(-2.5, 1, -1, 1, colors)
		imgf, err := os.Create("testdata/sample_img.dat") // Is there a cleaner way to do this/
		if err != nil {
			panic(err)
		}
		defer f.Close()
		gob.Register(&image.RGBA{})
		gob.Register(color.RGBA{})
		enc = gob.NewEncoder(imgf)
		err = enc.Encode(img)
		if err != nil {
			panic(err)
		}

	}
	m.Run()
}

// TestGif _wants_ to test creating a small mandelbrot gif, but can;t encode a sample
func TestGif(t *testing.T) {
	colors := mandelbrot.NewPalette(10)
	got := mandelbrot.MandelbrotBuilder{10, 10, 10}.Gif(3, -1.5, 0, .98, colors)
	want := &gif.GIF{}
	f, err := os.Open("testdata/sample_gif.dat") // Is there a cleaner way to do this/
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	gob.Register(&gif.GIF{})
	gob.Register(color.RGBA{})
	enc := gob.NewDecoder(f)
	err = enc.Decode(want)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

// TestDraw tests creating a small mandelbrot image and compares to a known-good image
func TestDraw(t *testing.T) {
	colors := mandelbrot.NewPalette(10)
	got := mandelbrot.MandelbrotBuilder{10, 5, 10}.Draw(-2.5, 1, -1, 1, colors)
	want := &image.RGBA{}
	f, err := os.Open("testdata/sample_img.dat")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	dec := gob.NewDecoder(f)
	err = dec.Decode(want)
	if err != nil {
		t.Errorf("decode error: %v", err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
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
		fn := mandelbrot.Scale(tc.inputMin, tc.inputMax, tc.outputMin, tc.outputMax)
		value := fn(tc.input)
		if value != tc.output {
			t.Fatalf("got %v for %v", value, tc.name)
		}
	}
}
