package main

import (
	"log"
	"net/http"
)

func main() {
	setupAdminServer()

	watchHostsFile()

	log.Fatal(http.ListenAndServe(":80", nil))
}
