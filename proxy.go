package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var proxy *httputil.ReverseProxy

func setupProxy() {
	origin, _ := url.Parse("http://127.0.0.1:4000/")
	proxy = httputil.NewSingleHostReverseProxy(origin)

	http.HandleFunc("/", handleRedirect)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s\n", r.Host, r.URL)
	proxy.ServeHTTP(w, r)
}
