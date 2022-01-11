package main

import (
	"embed"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"text/template"
)

type ProxyErrorTemplateData struct {
	Title     string
	Message   string
	ErrorCode string
}

//go:embed assets/*
var assets embed.FS

var proxies map[string]*httputil.ReverseProxy

func parseAndSetupProxies(hostsFilePath string) {
	setupProxies(NewHostsFile(hostsFilePath))
}

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
		w.WriteHeader(502)
		displayProxyError(w, r, &ProxyErrorTemplateData{
			Title:     "Proxy Entry Not Found",
			Message:   fmt.Sprintf("<p>Proxy entry not found for %s.</p>", r.Host),
			ErrorCode: "wdp1000",
		})
	}
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	errText := err.Error()

	var templateData *ProxyErrorTemplateData

	if strings.Contains(errText, "connection refused") {
		templateData = &ProxyErrorTemplateData{
			Title:     "Proxied Request Refused",
			Message:   fmt.Sprintf("<p>Proxied request to <a href=\"%s\">%s</a> refused.</p><p>Make sure the destination server is running on port %s.</p>", r.URL, r.URL, r.URL.Port()),
			ErrorCode: "wdp1001",
		}
	} else {
		templateData = &ProxyErrorTemplateData{
			Title:     "Proxy Error",
			Message:   fmt.Sprintf("<p>Proxy error: %s</p>", err.Error()),
			ErrorCode: "wdp2000",
		}
	}

	displayProxyError(w, r, templateData)
}

func displayProxyError(w http.ResponseWriter, r *http.Request, templateData *ProxyErrorTemplateData) {
	t, _ := template.ParseFS(assets, "assets/proxy_error.html")
	t.Execute(w, templateData)
}
