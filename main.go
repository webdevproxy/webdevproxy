package main

import (
	"log"
	"net/http"
)

const hostsFilePath = "/etc/hosts"

type HostWatcherSub struct{}

func (hw *HostWatcherSub) receive(path, event string) {
	parseAndSetupProxies(path)
}

func main() {
	parseAndSetupProxies(hostsFilePath)

	var hostsWatcher FSWatcherPub = &FSWatcher{
		path: hostsFilePath,
	}
	go hostsWatcher.observe()

	var hostsWatcherSub FSWatcherSub = &HostWatcherSub{}
	hostsWatcher.register(&hostsWatcherSub)

	log.Fatal(http.ListenAndServe(":80", nil))
}
