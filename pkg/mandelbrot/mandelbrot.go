package mandelbrot

import (
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/png"
	"math"
	"os"
	"sync"
)

// FloatFunction is a takes a float64 and returns a float64.
type FloatFunction func(a float64) float64

// Gif returns the gif containing frames and delays for a mandelbrot animation
func Gif(sizeX, sizeY, frames uint16, maxIterations uint8, x, y, scaleIn float64, colors []color.Color) *gif.GIF {
	var images []*image.Paletted
	var delays []int
	xShift := 1.0
	yShift := 1.0
	xMin, xMax, yMin, yMax := ExtentFromPoint(x, y, xShift, yShift)
	for frame := uint16(0); frame < frames; frame++ {
		img := Draw(sizeX, sizeY, maxIterations, xMin, xMax, yMin, yMax, colors)
		palettedImage := image.NewPaletted(img.Bounds(), colors)
		draw.Draw(palettedImage, palettedImage.Rect, img, img.Bounds().Min, draw.Over)
		images = append(images, palettedImage)
		delays = append(delays, 0)
		xShift *= scaleIn
		yShift *= scaleIn
		xMin, xMax, yMin, yMax = ExtentFromPoint(x, y, xShift, yShift)
	}
	return &gif.GIF{
		Image: images,
		Delay: delays,
	}
}

// ExtentFromPoint converts an x,y point + offsets in x and y into ranges of x and y values
func ExtentFromPoint(x, y, xShift, yShift float64) (xMin, xMax, yMin, yMax float64) {
	xMin = x - xShift
	xMax = x + xShift
	yMin = y - yShift
	yMax = y + yShift
	return
}

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
		colors = append(colors, color.RGBA{1, value, value, 1})
	}
	return
}

// WritePng writes an image to a filename
// Is there a good way to test this without deleting and recreating the file?
func WritePng(img *image.RGBA, filename string) {
	f, _ := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	png.Encode(f, img)
}

// EscapeIterations calculates how many iterations it takes for this point to escape Mandelbrot iteration, with a cap of maxIterations
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

// HasEscaped tells us whether a point has escaped under Mandelbrot iteration, ie it has length > 2
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
