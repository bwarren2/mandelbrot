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
	req, err := http.NewRequest("GET", "healthcheck", nil)
	check(err)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(web.HealthcheckHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status code!  Got %v wanted %v", status, http.StatusOK)
	}
	if rr.Body.String() != "heartbeat" {
		t.Errorf("Wrong body!  Got %v wanted %v", rr.Body.String(), "heartbeat")
	}
}
