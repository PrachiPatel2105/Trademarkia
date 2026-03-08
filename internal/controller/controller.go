package controller

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"hotreload/internal/builder"
	"hotreload/internal/cli"
	"hotreload/internal/debounce"
	"hotreload/internal/process"
	"hotreload/internal/watcher"
)

// Controller coordinates the hot reload workflow
type Controller struct {
	config  *cli.Config
	log     *slog.Logger
	watcher *watcher.Watcher
	builder *builder.Builder
	process *process.Manager

	mu            sync.Mutex
	buildPending  bool
	building      bool
	restartCount  int
	lastRestart   time.Time
	debouncer     *debounce.Debouncer
}

// New creates a new Controller
func New(config *cli.Config, log *slog.Logger) (*Controller, error) {
	ctrl := &Controller{
		config: config,
		log:    log,
	}

	// Create builder
	ctrl.builder = builder.New(config.BuildCmd, config.Root, log)

	// Create process manager
	ctrl.process = process.New(config.ExecCmd, config.Root, log)

	// Create debouncer (500ms delay)
	ctrl.debouncer = debounce.New(500*time.Millisecond, ctrl.onDebounced)

	// Create watcher
	w, err := watcher.New(config.Root, log, ctrl.onFileChange)
	if err != nil {
		return nil, err
	}
	ctrl.watcher = w

	return ctrl, nil
}

// Run starts the controller
func (ctrl *Controller) Run(ctx context.Context) error {
	// Start debouncer
	go ctrl.debouncer.Start(ctx)

	// Start watcher
	go func() {
		if err := ctrl.watcher.Start(ctx); err != nil {
			ctrl.log.Error("watcher error", "error", err)
		}
	}()

	// Initial build and start
	if err := ctrl.buildAndRestart(ctx); err != nil {
		ctrl.log.Error("initial build failed", "error", err)
		return err
	}

	// Wait for context cancellation
	<-ctx.Done()

	// Stop server
	if err := ctrl.process.Stop(); err != nil {
		ctrl.log.Error("failed to stop server", "error", err)
	}

	return nil
}

// onFileChange is called when a file changes
func (ctrl *Controller) onFileChange(path string) {
	ctrl.log.Debug("file change detected", "path", path)
	ctrl.debouncer.Trigger()
}

// onDebounced is called after debounce period
func (ctrl *Controller) onDebounced() {
	ctrl.mu.Lock()
	if ctrl.building {
		// Mark that another build is needed
		ctrl.buildPending = true
		ctrl.mu.Unlock()
		return
	}
	ctrl.building = true
	ctrl.mu.Unlock()

	// Build and restart in background
	go func() {
		ctx := context.Background()
		if err := ctrl.buildAndRestart(ctx); err != nil {
			ctrl.log.Error("rebuild failed", "error", err)
		}

		ctrl.mu.Lock()
		ctrl.building = false
		pending := ctrl.buildPending
		ctrl.buildPending = false
		ctrl.mu.Unlock()

		// If another build was requested, trigger it
		if pending {
			ctrl.onDebounced()
		}
	}()
}

// buildAndRestart builds the project and restarts the server
func (ctrl *Controller) buildAndRestart(ctx context.Context) error {
	// Check for restart loop protection
	if ctrl.shouldThrottle() {
		ctrl.log.Warn("too many restarts, throttling for 5 seconds")
		time.Sleep(5 * time.Second)
	}

	// Build
	if err := ctrl.builder.Build(ctx); err != nil {
		return err
	}

	// Stop existing process
	if ctrl.process.IsRunning() {
		if err := ctrl.process.Stop(); err != nil {
			ctrl.log.Error("failed to stop process", "error", err)
		}
	}

	// Start new process
	if err := ctrl.process.Start(ctx); err != nil {
		return err
	}

	// Track restart
	ctrl.mu.Lock()
	ctrl.restartCount++
	ctrl.lastRestart = time.Now()
	ctrl.mu.Unlock()

	return nil
}

// shouldThrottle checks if we should throttle restarts
func (ctrl *Controller) shouldThrottle() bool {
	ctrl.mu.Lock()
	defer ctrl.mu.Unlock()

	// If more than 5 restarts in the last 10 seconds, throttle
	if ctrl.restartCount > 5 && time.Since(ctrl.lastRestart) < 10*time.Second {
		ctrl.restartCount = 0
		return true
	}

	// Reset counter if enough time has passed
	if time.Since(ctrl.lastRestart) > 10*time.Second {
		ctrl.restartCount = 0
	}

	return false
}
