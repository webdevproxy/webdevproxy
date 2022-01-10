package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

var proxies map[string]*httputil.ReverseProxy

func setupProxies(hostsfile *Hostfile) {
	haveExistingProxies := proxies != nil
	proxies = make(map[string]*httputil.ReverseProxy)

	for _, entry := range hostsfile.Entries {
		if entry.Proxied {
			proxy := httputil.NewSingleHostReverseProxy(entry.proxyUrl())
			proxy.ErrorHandler = handleError
			proxies[entry.Hostname] = proxy
		}
	}

	if !haveExistingProxies {
		http.HandleFunc("/", handleRedirect)
	}
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	if proxies[r.Host] != nil {
		proxies[r.Host].ServeHTTP(w, r)
	} else {
		w.WriteHeader(500)
		fmt.Fprint(w, "Proxy not found for request")
	}
}

func handleError(res http.ResponseWriter, req *http.Request, err error) {
	fmt.Fprintf(res, "Connection refused when proxying to %s", req.URL)
}
