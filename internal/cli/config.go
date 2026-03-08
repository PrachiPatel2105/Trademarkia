package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the CLI configuration
type Config struct {
	Root     string
	BuildCmd string
	ExecCmd  string
}

// ParseFlags parses command-line flags and returns configuration
func ParseFlags() *Config {
	config := &Config{}

	flag.StringVar(&config.Root, "root", ".", "Project root directory to watch")
	flag.StringVar(&config.BuildCmd, "build", "", "Build command (required)")
	flag.StringVar(&config.ExecCmd, "exec", "", "Execute command (required)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: hotreload [options]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  hotreload --root ./myproject --build \"go build -o ./bin/server ./cmd/server\" --exec \"./bin/server\"\n")
	}

	flag.Parse()

	return config
}

// Validate checks if the configuration is valid
func Validate(config *Config) error {
	if config.BuildCmd == "" {
		return errors.New("--build flag is required")
	}

	if config.ExecCmd == "" {
		return errors.New("--exec flag is required")
	}

	// Check if root directory exists
	absRoot, err := filepath.Abs(config.Root)
	if err != nil {
		return fmt.Errorf("invalid root path: %w", err)
	}

	info, err := os.Stat(absRoot)
	if err != nil {
		return fmt.Errorf("root directory does not exist: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("root path is not a directory: %s", absRoot)
	}

	config.Root = absRoot
	return nil
}
