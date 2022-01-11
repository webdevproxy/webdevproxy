package main

const hostsFilePath = "/etc/hosts"

var hostsWatcher FSWatcherPub

type HostsWatcherSub struct{}

func (hw *HostsWatcherSub) receive(path, event string) {
	parseAndSetupProxies(path)
}

func watchHostsFile() {
	hostsWatcher = &FSWatcher{
		path: hostsFilePath,
	}
	go hostsWatcher.observe()

	var hostsWatcherSub FSWatcherSub = &HostsWatcherSub{}
	hostsWatcher.register(&hostsWatcherSub)

	parseAndSetupProxies(hostsFilePath)
}
