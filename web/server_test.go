package web_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bwarren2/mandelbrot/web"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
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
		{"GET", "mandelbrot/png", web.MandelPngHandler, 200, "",
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
		handler := http.HandlerFunc(tc.handler)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != tc.responseCode {
			t.Errorf("Wrong status code!  Got %v wanted %v", status, tc.responseCode)
		}
		if tc.responseBody != "" && rr.Body.String() != tc.responseBody {
			t.Errorf("Wrong body!  Got %v wanted %v", rr.Body.String(), tc.responseBody)
		}
	}
}
