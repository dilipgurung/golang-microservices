package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.HandleFunc("/healthcheck", healthCheckHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func healthCheckHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
}
