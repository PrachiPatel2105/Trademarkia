package builder

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
	"time"
)

// Builder executes build commands
type Builder struct {
	buildCmd string
	workDir  string
	log      *slog.Logger
}

// New creates a new Builder
func New(buildCmd, workDir string, log *slog.Logger) *Builder {
	return &Builder{
		buildCmd: buildCmd,
		workDir:  workDir,
		log:      log,
	}
}

// Build executes the build command
func (b *Builder) Build(ctx context.Context) error {
	b.log.Info("build started", "command", b.buildCmd)
	start := time.Now()

	// Parse command
	parts := parseCommand(b.buildCmd)
	if len(parts) == 0 {
		return fmt.Errorf("empty build command")
	}

	// Create command
	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
	cmd.Dir = b.workDir

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Execute
	err := cmd.Run()
	duration := time.Since(start)

	if err != nil {
		b.log.Error("build failed",
			"duration", duration,
			"error", err)

		// Print build errors
		if stderr.Len() > 0 {
			fmt.Print("\n--- Build Errors ---\n")
			fmt.Print(stderr.String())
			fmt.Print("--- End Build Errors ---\n\n")
		}

		return fmt.Errorf("build failed: %w", err)
	}

	b.log.Info("build succeeded", "duration", duration)

	// Print build output if any
	if stdout.Len() > 0 {
		fmt.Print(stdout.String())
	}

	return nil
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
