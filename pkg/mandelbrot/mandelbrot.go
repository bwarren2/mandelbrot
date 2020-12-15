package mandelbrot

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"sync"
)

// Retro is a color scheme borrowed from the internet
var Retro = []color.Color{
	color.RGBA{0x00, 0x04, 0x0f, 0xff},
	color.RGBA{0x03, 0x26, 0x28, 0xff},
	color.RGBA{0x07, 0x3e, 0x1e, 0xff},
	color.RGBA{0x18, 0x55, 0x08, 0xff},
	color.RGBA{0x5f, 0x6e, 0x0f, 0xff},
	color.RGBA{0x84, 0x50, 0x19, 0xff},
	color.RGBA{0x9b, 0x30, 0x22, 0xff},
	color.RGBA{0xb4, 0x92, 0x2f, 0xff},
	color.RGBA{0x94, 0xca, 0x3d, 0xff},
	color.RGBA{0x4f, 0xd5, 0x51, 0xff},
	color.RGBA{0x66, 0xff, 0xb3, 0xff},
	color.RGBA{0x82, 0xc9, 0xe5, 0xff},
	color.RGBA{0x9d, 0xa3, 0xeb, 0xff},
	color.RGBA{0xd7, 0xb5, 0xf3, 0xff},
	color.RGBA{0xfd, 0xd6, 0xf6, 0xff},
	color.RGBA{0xff, 0xf0, 0xf2, 0xff},
}

// FloatFunction is a takes a float64 and returns a float64.
type FloatFunction func(a float64) float64

// ColorRow fills in one row of Mandelbrot values for an image
func ColorRow(img *image.RGBA, row, length uint16, xScale, yScale FloatFunction, colors []color.Color, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := uint16(0); j < length; j++ {
		pointX := xScale(float64(row))
		pointY := yScale(float64(j))
		iterations := EscapeIterations(pointX, pointY, 300)
		color := colors[iterations%len(colors)]
		img.Set(int(row), int(j), color)
	}
}

// Draw draws a Mandelbrot image of a given size with a given domain and range
func Draw(sizeX, sizeY uint16, maxIterations uint8, minX, maxX, minY, maxY float64, colors []color.Color) *image.RGBA {
	var wg sync.WaitGroup
	img := image.NewRGBA(image.Rect(0, 0, int(sizeX), int(sizeY)))
	xScale := Scale(0, float64(sizeX), float64(minX), float64(maxX))
	yScale := Scale(0, float64(sizeY), float64(minY), float64(maxY))

	for i := uint16(0); i < sizeX; i++ {
		wg.Add(1)
		go ColorRow(img, i, sizeY, xScale, yScale, colors, &wg)
	}
	wg.Wait()
	return img
}

// NewPalette returns a list of colors to use as a palette
func NewPalette(maxIterations uint8) (colors []color.Color) {
	colorScale := Scale(0, float64(maxIterations), 0, 255)
	for x := uint8(0); x < maxIterations; x++ {
		value := uint8(colorScale(float64(x)))
		colors = append(colors, color.RGBA{value, value, value, 1})
	}
	return
}

// WritePng writes an image to a filename
func WritePng(img *image.RGBA, filename string) {
	f, _ := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	png.Encode(f, img)
}

// EscapeIterations calculates how many iterations it takes for this point to escape Mandelbrot application
func EscapeIterations(x0, y0 float64, maxIterations int) (iterations int) {
	x := 0.0
	y := 0.0
	var xTemp float64
	for !HasEscaped(x, y) && iterations < maxIterations {
		xTemp = math.Pow(x, 2) - math.Pow(y, 2) + x0
		y = 2*x*y + y0
		x = xTemp
		iterations++
	}
	return
}

// HasEscaped tells us whether a point has escaped under Mandelbrot iteration
func HasEscaped(x, y float64) bool {
	return math.Pow(x, 2)+math.Pow(y, 2) > 4
}

// Scale returns a scaling function clamped to a given range
func Scale(inputMin, inputMax, outputMin, outputMax float64) func(a float64) float64 {
	return func(a float64) float64 {
		if a < math.Min(inputMin, inputMax) {
			return outputMin
		} else if a > math.Max(inputMin, inputMax) {
			return outputMax
		}
		return outputMin + (outputMax-outputMin)*(a-inputMin)/(inputMax-inputMin)
	}
}
