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

// Draw draws a Mandelbrot image of a given size with a given domain and range
func Draw(filename string, sizeX, sizeY int, minX, maxX, minY, maxY float64) {
	var wg sync.WaitGroup
	img := image.NewRGBA(image.Rect(0, 0, sizeX, sizeY))
	xScale := Scale(0, float64(sizeX), float64(minX), float64(maxX))
	yScale := Scale(0, float64(sizeY), float64(minY), float64(maxY))
	colorScale := Scale(20, 0, 0, 255)
	for i := 0; i < sizeX; i++ {
		wg.Add(1)
		go func(img *image.RGBA, row, length int, xScale, yScale, colorScale func(a float64) float64, wg *sync.WaitGroup) {
			defer wg.Done()

			for j := 0; j < sizeY; j++ {
				pointX := xScale(float64(row))
				pointY := yScale(float64(j))
				iterations := EscapeIterations(pointX, pointY, 300)
				color := Retro[iterations%len(Retro)]
				img.Set(row, j, color)
			}
		}(img, i, sizeY, xScale, yScale, colorScale, &wg)
	}
	wg.Wait()
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
