package web

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/bwarren2/mandelbrot"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// DefaultMap is a string map that supports get-with-default operations
type DefaultMap map[string]string

// GetDefaultInt gets an int from the DefaultMap if it exists and is an int, otherwise returning the default
func (d DefaultMap) GetDefaultInt(key string, otherwise int) int {
	if v, ok := d[key]; ok {
		value, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return otherwise
		}
		return int(value)
	}
	return otherwise
}

// GetDefaultFloat64 gets a float from the DefaultMap if it exists and is a float, otherwise returning the default
func (d DefaultMap) GetDefaultFloat64(key string, otherwise float64) float64 {
	if v, ok := d[key]; ok {
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return otherwise
		}
		return value
	}
	return otherwise
}

// NewDefaultMap makes a map that supports get-with-default
func NewDefaultMap(input url.Values) DefaultMap {
	newMap := make(DefaultMap)
	for key, value := range input {
		newMap[key] = value[0]
	}
	return newMap
}

// HealthcheckHandler just prints "heartbeat" to confirm the server is running
func HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "heartbeat")
}

// MandelPngHandlerGenerator takes a Drawer and returns a HttpRequestHandler.  This pattern supports dependency injection for testing.
func MandelPngHandlerGenerator(drawer mandelbrot.Drawer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		colors := mandelbrot.NewPalette(10)
		queryMap := NewDefaultMap(r.URL.Query())
		drawer.SetSize(uint16(queryMap.GetDefaultInt("sizeX", 1000)), uint16(queryMap.GetDefaultInt("sizeY", 500)))
		drawer.SetIterations(uint8(queryMap.GetDefaultInt("maxIterations", 10)))
		img := drawer.Draw(
			queryMap.GetDefaultFloat64("minX", -2.5),
			queryMap.GetDefaultFloat64("maxX", 1),
			queryMap.GetDefaultFloat64("minY", -1),
			queryMap.GetDefaultFloat64("maxY", 1),
			colors,
		)
		WriteImage(w, img)
	}
}

// WriteImage writes an image to a ResponseWriter and sets headers
func WriteImage(w http.ResponseWriter, img *image.RGBA) {
	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, img, nil); err != nil {
		log.Fatal("Could not write")
	}
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Fatal("Could not write")
	}
}

// Serve starts a server that supports healthcheck and a mandelbrot-png endpoint
func Serve(port int) {
	http.HandleFunc("/healthcheck", HealthcheckHandler)
	http.HandleFunc("/mandelbrot/png", MandelPngHandlerGenerator(&mandelbrot.MandelbrotBuilder{}))
	log.Fatal(http.ListenAndServe(":"+fmt.Sprint(port), nil))
}
