package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AdminPongInfo struct {
	Version string `json:"version"`
}

type AdminPong struct {
	Webdevproxy AdminPongInfo `json:"webdevproxy"`
}

type AdminWatchConfigMessage struct {
	Hosts *Hostsfile `json:"hosts"`
}

var adminMux *http.ServeMux
var adminStaticHandler http.Handler

func setupAdminServer() {
	adminMux = http.NewServeMux()
	adminMux.Handle("/", http.FileServer(http.FS(adminContent())))
	adminMux.HandleFunc("/ping", handleAdminServerPingRequest)
	adminMux.HandleFunc("/api/watch-config", handleAdminServerWatchConfigRequest)
}

func handleAdminServerRequest(w http.ResponseWriter, r *http.Request) {
	adminMux.ServeHTTP(w, r)
}

func handleAdminServerPingRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "webdevproxy.com")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type")
	w.Header().Set("Allow", "POST, OPTIONS")

	switch r.Method {
	case http.MethodPost:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AdminPong{
			Webdevproxy: AdminPongInfo{
				Version: "1.0.0",
			},
		})

	case http.MethodOptions:
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, fmt.Sprintf("%s method not allowed for %s", r.Method, r.URL.Path), http.StatusMethodNotAllowed)
	}
}

func handleAdminServerWatchConfigRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	/*
		TODO:

		1. convert hosts watcher to receive channel instead of callback:
			https://eli.thegreenplace.net/2020/pubsub-using-channels-in-go/ (look at bottom!)
		2. auto send current value when joining watcher
	*/

	hostsfileChan := make(chan *Hostsfile)

	go func() {
		hostsfileChan <- NewHostsFile(hostsFilePath)
	}()

	defer func() {
		close(hostsfileChan)
		hostsfileChan = nil
	}()

	flusher, _ := w.(http.Flusher)

	for {
		select {

		case hostfile := <-hostsfileChan:
			fmt.Fprint(w, "data: ")
			json.NewEncoder(w).Encode(AdminWatchConfigMessage{
				Hosts: hostfile,
			})
			fmt.Fprint(w, "\n\n")
			flusher.Flush()

		case <-r.Context().Done():
			return
		}
	}
}
