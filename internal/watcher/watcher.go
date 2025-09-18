package watcher

import "github.com/fsnotify/fsnotify"

type FileWatcher struct {
	watcher *fsnotify.Watcher
}

// TODO: implement general filewatcher interface
