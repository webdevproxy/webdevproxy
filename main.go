package main

import (
	"log"
	"net/http"
)

func main() {
	watchHostsFile()

	log.Fatal(http.ListenAndServe(":80", nil))
}
