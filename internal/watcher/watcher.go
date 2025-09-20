package watcher

import "github.com/fsnotify/fsnotify"

type Watcher struct {
	watcher *fsnotify.Watcher
}

// TODO: implement general filewatcher interface
