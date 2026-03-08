package watcher

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"

	"hotreload/internal/filter"
)

// Watcher watches file system changes
type Watcher struct {
	root     string
	watcher  *fsnotify.Watcher
	log      *slog.Logger
	onChange func(path string)
}

// New creates a new Watcher
func New(root string, log *slog.Logger, onChange func(path string)) (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	w := &Watcher{
		root:     root,
		watcher:  fsWatcher,
		log:      log,
		onChange: onChange,
	}

	// Add all directories recursively
	if err := w.addRecursive(root); err != nil {
		fsWatcher.Close()
		return nil, err
	}

	return w, nil
}

// addRecursive adds all directories under path to the watcher
func (w *Watcher) addRecursive(path string) error {
	return filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip ignored paths
		if filter.ShouldIgnore(p) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Only watch directories
		if info.IsDir() {
			if err := w.watcher.Add(p); err != nil {
				w.log.Warn("failed to watch directory", "path", p, "error", err)
			} else {
				w.log.Debug("watching directory", "path", p)
			}
		}

		return nil
	})
}

// Start begins watching for file changes
func (w *Watcher) Start(ctx context.Context) error {
	w.log.Info("watching directory", "root", w.root)

	for {
		select {
		case <-ctx.Done():
			return w.watcher.Close()

		case event, ok := <-w.watcher.Events:
			if !ok {
				return nil
			}

			w.handleEvent(event)

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return nil
			}
			w.log.Error("watcher error", "error", err)
		}
	}
}

// handleEvent processes a file system event
func (w *Watcher) handleEvent(event fsnotify.Event) {
	path := event.Name

	// Ignore filtered paths
	if filter.ShouldIgnore(path) {
		return
	}

	// Handle different event types
	switch {
	case event.Op&fsnotify.Create == fsnotify.Create:
		w.handleCreate(path)
	case event.Op&fsnotify.Write == fsnotify.Write:
		w.handleWrite(path)
	case event.Op&fsnotify.Remove == fsnotify.Remove:
		w.handleRemove(path)
	case event.Op&fsnotify.Rename == fsnotify.Rename:
		w.handleRemove(path)
	}
}

// handleCreate handles file/directory creation
func (w *Watcher) handleCreate(path string) {
	info, err := os.Stat(path)
	if err != nil {
		return
	}

	if info.IsDir() {
		// Add new directory to watcher
		if err := w.addRecursive(path); err != nil {
			w.log.Warn("failed to watch new directory", "path", path, "error", err)
		} else {
			w.log.Info("watching new directory", "path", path)
		}
	} else {
		// File created
		w.log.Info("file created", "path", path)
		w.onChange(path)
	}
}

// handleWrite handles file modifications
func (w *Watcher) handleWrite(path string) {
	info, err := os.Stat(path)
	if err != nil {
		return
	}

	if !info.IsDir() {
		w.log.Info("file changed", "path", path)
		w.onChange(path)
	}
}

// handleRemove handles file/directory removal
func (w *Watcher) handleRemove(path string) {
	w.log.Info("file removed", "path", path)
	// fsnotify automatically removes watches for deleted paths
	w.onChange(path)
}
