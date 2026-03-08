package process

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"
)

// Manager manages the server process lifecycle
type Manager struct {
	execCmd string
	workDir string
	log     *slog.Logger

	mu      sync.Mutex
	cmd     *exec.Cmd
	cancel  context.CancelFunc
	running bool
}

// New creates a new process Manager
func New(execCmd, workDir string, log *slog.Logger) *Manager {
	return &Manager{
		execCmd: execCmd,
		workDir: workDir,
		log:     log,
	}
}

// Start starts the server process
func (m *Manager) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		return fmt.Errorf("process already running")
	}

	// Parse command
	parts := parseCommand(m.execCmd)
	if len(parts) == 0 {
		return fmt.Errorf("empty exec command")
	}

	// Create cancellable context
	procCtx, cancel := context.WithCancel(ctx)
	m.cancel = cancel

	// Create command with process group
	cmd := exec.CommandContext(procCtx, parts[0], parts[1:]...)
	cmd.Dir = m.workDir

	// Set process group ID for killing child processes
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}

	// Setup stdout/stderr streaming
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		cancel()
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start process
	if err := cmd.Start(); err != nil {
		cancel()
		return fmt.Errorf("failed to start process: %w", err)
	}

	m.cmd = cmd
	m.running = true

	m.log.Info("server started", "pid", cmd.Process.Pid, "command", m.execCmd)

	// Stream output in goroutines
	go m.streamOutput(stdout, "stdout")
	go m.streamOutput(stderr, "stderr")

	// Monitor process
	go m.monitor()

	return nil
}

// Stop stops the server process
func (m *Manager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running {
		return nil
	}

	m.log.Info("stopping server", "pid", m.cmd.Process.Pid)

	// Cancel context
	if m.cancel != nil {
		m.cancel()
	}

	// Try graceful shutdown first (SIGTERM equivalent on Windows)
	if err := m.cmd.Process.Signal(os.Interrupt); err != nil {
		m.log.Warn("failed to send interrupt signal", "error", err)
	}

	// Wait for graceful shutdown with timeout
	done := make(chan error, 1)
	go func() {
		done <- m.cmd.Wait()
	}()

	select {
	case <-time.After(5 * time.Second):
		// Force kill if not stopped
		m.log.Warn("process did not stop gracefully, forcing kill")
		if err := m.cmd.Process.Kill(); err != nil {
			m.log.Error("failed to kill process", "error", err)
		}
		<-done // Wait for process to be reaped
	case <-done:
		// Process stopped
	}

	m.running = false
	m.cmd = nil
	m.cancel = nil

	m.log.Info("server stopped")
	return nil
}

// IsRunning returns true if the process is running
func (m *Manager) IsRunning() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.running
}

// streamOutput streams process output to stdout/stderr
func (m *Manager) streamOutput(reader io.Reader, name string) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if name == "stderr" {
			fmt.Fprintf(os.Stderr, "[server] %s\n", line)
		} else {
			fmt.Printf("[server] %s\n", line)
		}
	}
}

// monitor watches the process and logs when it exits
func (m *Manager) monitor() {
	err := m.cmd.Wait()

	m.mu.Lock()
	m.running = false
	m.mu.Unlock()

	if err != nil {
		m.log.Warn("server exited with error", "error", err)
	} else {
		m.log.Info("server exited")
	}
}

// parseCommand splits a command string into parts, respecting quotes
func parseCommand(cmd string) []string {
	var parts []string
	var current strings.Builder
	inQuote := false
	quoteChar := rune(0)

	for _, r := range cmd {
		switch {
		case r == '"' || r == '\'':
			if inQuote {
				if r == quoteChar {
					inQuote = false
					quoteChar = 0
				} else {
					current.WriteRune(r)
				}
			} else {
				inQuote = true
				quoteChar = r
			}
		case r == ' ' && !inQuote:
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}
