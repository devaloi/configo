package watcher

import (
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watcher monitors a file for changes and calls subscribers after debounce.
type Watcher struct {
	path      string
	debounce  time.Duration
	onChange  []func()
	fsWatcher *fsnotify.Watcher
	done      chan struct{}
	mu        sync.Mutex
}

// New creates a Watcher for the given file path.
func New(path string, debounce time.Duration) *Watcher {
	if debounce == 0 {
		debounce = 500 * time.Millisecond
	}
	return &Watcher{
		path:     path,
		debounce: debounce,
		done:     make(chan struct{}),
	}
}

// OnChange registers a callback that fires when the watched file changes.
func (w *Watcher) OnChange(fn func()) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.onChange = append(w.onChange, fn)
}

// Start begins watching the file. It blocks until Stop is called.
func (w *Watcher) Start() error {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	w.fsWatcher = fw

	if err := fw.Add(w.path); err != nil {
		fw.Close()
		return err
	}

	go w.loop()
	return nil
}

func (w *Watcher) loop() {
	var timer *time.Timer
	for {
		select {
		case event, ok := <-w.fsWatcher.Events:
			if !ok {
				return
			}
			if event.Op&(fsnotify.Write|fsnotify.Create) == 0 {
				continue
			}
			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(w.debounce, func() {
				w.mu.Lock()
				handlers := make([]func(), len(w.onChange))
				copy(handlers, w.onChange)
				w.mu.Unlock()
				for _, fn := range handlers {
					fn()
				}
			})
		case _, ok := <-w.fsWatcher.Errors:
			if !ok {
				return
			}
		case <-w.done:
			if timer != nil {
				timer.Stop()
			}
			return
		}
	}
}

// Stop stops watching the file.
func (w *Watcher) Stop() error {
	close(w.done)
	if w.fsWatcher != nil {
		return w.fsWatcher.Close()
	}
	return nil
}
