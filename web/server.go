package web

import (
	"fmt"
	"log"
	"net/http"
)

func HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "heartbeat")
}

func Serve(port int) {
	http.HandleFunc("/healthcheck", HealthcheckHandler)
	log.Fatal(http.ListenAndServe(":"+fmt.Sprint(port), nil))
}
