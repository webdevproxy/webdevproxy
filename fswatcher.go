// adapted from https://github.com/htmfilho/blog-examples (CC0 licence)

package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type FSWatcherPub interface {
	register(subscriber *FSWatcherSub)
	unregister(subscriber *FSWatcherSub)
	notify(path, event string)
	observe()
}

type FSWatcherSub interface {
	receive(path, event string)
}

type FSWatcher struct {
	subscribers []*FSWatcherSub
	watcher     fsnotify.Watcher
	path        string
}

func (pw *FSWatcher) register(subscriber *FSWatcherSub) {
	pw.subscribers = append(pw.subscribers, subscriber)
}

func (pw *FSWatcher) unregister(subscriber *FSWatcherSub) {
	length := len(pw.subscribers)

	for i, sub := range pw.subscribers {
		if sub == subscriber {
			pw.subscribers[i] = pw.subscribers[length-1]
			pw.subscribers = pw.subscribers[:length-1]
			break
		}
	}
}

func (pw *FSWatcher) notify(path, event string) {
	for _, sub := range pw.subscribers {
		(*sub).receive(path, event)
	}
}

func (pw *FSWatcher) observe() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Error", err)
	}
	defer watcher.Close()

	fileInfo, err := os.Stat(pw.path)
	if err == nil && fileInfo.IsDir() {
		filepath.WalkDir(pw.path, func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				return watcher.Add(path)
			}
			return nil
		})
	} else {
		watcher.Add(pw.path)
	}

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				pw.notify(event.Name, event.Op.String())
			case err := <-watcher.Errors:
				fmt.Println("Error", err)
			}
		}
	}()

	<-done
}
