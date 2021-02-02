package web_test

import (
	"image"
	"image/color"
	"image/gif"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bwarren2/mandelbrot"
	"github.com/bwarren2/mandelbrot/web"
	"github.com/google/go-cmp/cmp"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type DrawSpy struct {
	callArgs []interface{}
}

func (ds *DrawSpy) Draw(sizeX, sizeY uint16, maxIterations uint8, minX, maxX, minY, maxY float64, colors []color.Color) *image.RGBA {
	ds.callArgs = append(ds.callArgs, sizeX, sizeY, maxIterations, minX, maxX, minY, maxY, colors)
	return &image.RGBA{}
}
func (ds *DrawSpy) Gif(sizeX, sizeY, frames uint16, maxIterations uint8, x, y, scaleIn float64, colors []color.Color) *gif.GIF {
	ds.callArgs = append(ds.callArgs, sizeX, sizeY, frames, maxIterations, x, y, scaleIn, colors)
	return &gif.GIF{}
}

func TestWebServing(t *testing.T) {
	tcs := []struct {
		verb, url    string
		handler      http.HandlerFunc
		responseCode int
		responseBody string
		params       map[string]string
	}{
		{"GET", "healthcheck", web.HealthcheckHandler, 200, "heartbeat", map[string]string{}},
	}
	for _, tc := range tcs {
		req, err := http.NewRequest(tc.verb, tc.url, nil)
		check(err)
		q := req.URL.Query()
		for k, v := range tc.params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tc.handler)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != tc.responseCode {
			t.Errorf("Wrong status code!  Got %v wanted %v", status, tc.responseCode)
		}
		if rr.Body.String() != tc.responseBody {
			t.Errorf("Wrong body!  Got %v wanted %v", rr.Body.String(), tc.responseBody)
		}
	}
}
func TestWebHandler(t *testing.T) {
	tcs := []struct {
		verb, url    string
		generator    func(drawer mandelbrot.Drawer) func(w http.ResponseWriter, r *http.Request)
		responseCode int
		params       map[string]string
	}{
		{"GET", "mandelbrot/png", web.MandelPngHandlerGenerator, 200,
			map[string]string{"sizeX": "10", "sizeY": "5", "maxIterations": "10", "minX": "-2.5", "maxX": "1", "minY": "-1", "maxY": "1"}},
	}
	for _, tc := range tcs {
		req, err := http.NewRequest(tc.verb, tc.url, nil)
		check(err)
		q := req.URL.Query()
		for k, v := range tc.params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
		rr := httptest.NewRecorder()
		spy := &DrawSpy{}
		handler := http.HandlerFunc(web.MandelPngHandlerGenerator(spy))
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != tc.responseCode {
			t.Errorf("Wrong status code!  Got %v wanted %v", status, tc.responseCode)
		}
		if diff := cmp.Diff(spy.callArgs, []interface{}{10, 5, 10, -2.5, 1, -1, 1}); diff != "" {
			t.Errorf("Call args mismatch (-want +got):\n%s", diff)
		}
	}
}
