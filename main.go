package main

import (
	"log"
	"net/http"
)

func main() {
	setupProxy()

	log.Fatal(http.ListenAndServe(":80", nil))
}
