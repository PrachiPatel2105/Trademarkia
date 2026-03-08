package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"hotreload/internal/cli"
	"hotreload/internal/controller"
	"hotreload/internal/logger"
)

func main() {
	// Parse CLI flags
	config := cli.ParseFlags()

	// Validate configuration
	if err := cli.Validate(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	// Initialize logger
	log := logger.New()

	log.Info("hotreload starting",
		"root", config.Root,
		"build", config.BuildCmd,
		"exec", config.ExecCmd)

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Info("received signal, shutting down", "signal", sig)
		cancel()
	}()

	// Create and run controller
	ctrl, err := controller.New(config, log)
	if err != nil {
		log.Error("failed to create controller", "error", err)
		os.Exit(1)
	}

	if err := ctrl.Run(ctx); err != nil {
		log.Error("controller error", "error", err)
		os.Exit(1)
	}

	log.Info("hotreload stopped")
}
