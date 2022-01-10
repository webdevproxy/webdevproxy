package main

import (
	"log"
	"net/http"
)

func main() {
	hostsfile, err := NewHostsFile("/etc/hosts")
	if err != nil {
		log.Fatal(err)
	}

	setupProxies(hostsfile)

	log.Fatal(http.ListenAndServe(":80", nil))
}
