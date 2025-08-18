package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	registerMux(mux)

	log.Print("starting server on http://localhost:4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
