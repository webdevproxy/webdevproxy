package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
)

type ProxyErrorTemplateData struct {
	Title   string
	Message string
}

var mutex sync.RWMutex
var proxies map[string]*httputil.ReverseProxy
var hostFileEntries map[string]*HostsFileEntry

func parseAndSetupProxies(hostsFilePath string) {
	setupProxies(NewHostsFile(hostsFilePath))
}

func setupProxies(hostsfile *Hostsfile) {
	mutex.Lock()

	haveExistingProxies := proxies != nil
	proxies = make(map[string]*httputil.ReverseProxy)
	hostFileEntries = make(map[string]*HostsFileEntry)

	for _, entry := range hostsfile.Entries {
		hostFileEntries[entry.Host] = &entry
		if entry.Proxied {
			proxy := httputil.NewSingleHostReverseProxy(entry.proxyUrl())
			proxy.ErrorHandler = handleError
			proxies[entry.Host] = proxy
		}
	}

	mutex.Unlock()

	if !haveExistingProxies {
		http.HandleFunc("/", handleProxy)
	}
}

func handleProxy(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Host, r.URL)

	if r.Host == "local-admin.webdevproxy.com" {
		handleAdminServerRequest(w, r)
		return
	}

	mutex.RLock()
	proxy := proxies[r.Host]
	hostFileEntry := hostFileEntries[r.Host]
	mutex.RUnlock()

	if proxy != nil {
		proxy.ServeHTTP(w, r)
		return
	}

	if hostFileEntry != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		displayProxyError(w, r, &ProxyErrorTemplateData{
			Title:   "Host Not Proxied",
			Message: fmt.Sprintf("<p>The requested host (<strong>%s</strong>) is in your local hosts file but does not have a webdevproxy comment.</p>", r.Host),
		})
		return
	}

	handlePassThrough(w, r)
}

// TODO: use https://github.com/hashicorp/golang-lru to cache proxies
func handlePassThrough(w http.ResponseWriter, r *http.Request) {
	remoteUrl, _ := url.Parse(fmt.Sprintf("%s://%s", r.URL.Scheme, r.URL.Host))
	proxy := httputil.NewSingleHostReverseProxy(remoteUrl)
	proxy.ErrorHandler = handleError
	proxy.ServeHTTP(w, r)
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	errText := err.Error()

	var templateData *ProxyErrorTemplateData

	if strings.Contains(errText, "connection refused") {
		templateData = &ProxyErrorTemplateData{
			Title:   "Proxied Request Refused",
			Message: fmt.Sprintf("<p>Proxied request to <a href=\"%s\">%s</a> refused.</p><p>Make sure the destination server is running on port %s.</p>", r.URL, r.URL, r.URL.Port()),
		}
	} else {
		templateData = &ProxyErrorTemplateData{
			Title:   "Proxy Error",
			Message: fmt.Sprintf("<p><em>%s</em></p>", err.Error()),
		}
	}

	displayProxyError(w, r, templateData)
}

func displayProxyError(w http.ResponseWriter, r *http.Request, templateData *ProxyErrorTemplateData) {
	proxyErrorTemplate().Execute(w, templateData)
}
