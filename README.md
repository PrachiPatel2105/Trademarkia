# hotreload

A production-quality CLI tool for automatically rebuilding and restarting Go servers when code changes are detected.

## Features

- **Automatic rebuild and restart** on file changes
- **Debounced file events** to handle editor quirks (300-500ms window)
- **Recursive directory watching** with dynamic folder detection
- **Smart filtering** of unnecessary files (.git, node_modules, tmp files, etc.)
- **Real-time log streaming** from the running server
- **Graceful process management** with SIGTERM → SIGKILL fallback
- **Restart loop protection** to prevent rapid crash cycles
- **Build error reporting** without starting the server on failure
- **Long-running stability** designed for hours of continuous operation

## Architecture

```
hotreload/
├── cmd/hotreload/          # CLI entry point
│   └── main.go
├── internal/
│   ├── cli/                # Command-line parsing and validation
│   │   └── config.go
│   ├── watcher/            # File system monitoring (fsnotify)
│   │   └── watcher.go
│   ├── debounce/           # Event aggregation
│   │   └── debounce.go
│   ├── filter/             # Path filtering logic
│   │   └── filter.go
│   ├── builder/            # Build command execution
│   │   └── builder.go
│   ├── process/            # Process lifecycle management
│   │   └── manager.go
│   ├── controller/         # Main orchestration logic
│   │   └── controller.go
│   └── logger/             # Structured logging (slog)
│       └── logger.go
├── testserver/             # Demo HTTP server
│   └── main.go
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Component Responsibilities

### CLI Module
- Parses command-line flags
- Validates arguments
- Provides usage information

### Watcher Module
- Recursively watches directories using fsnotify
- Handles file system events (create, write, remove, rename)
- Dynamically adds new directories to watch list
- Filters ignored paths

### Debounce Module
- Aggregates rapid file events into single rebuild trigger
- Prevents multiple rebuilds for single save operation
- Configurable delay window (default: 500ms)

### Filter Module
- Ignores unnecessary directories (.git, node_modules, bin, build, dist, tmp)
- Filters temporary files (*.swp, *.tmp, *.log)
- Skips hidden files and editor artifacts

### Builder Module
- Executes build commands
- Captures and displays build output
- Reports build failures without starting server

### Process Manager Module
- Starts and stops server processes
- Streams stdout/stderr in real-time
- Implements graceful shutdown (SIGTERM → SIGKILL)
- Kills child processes via process groups
- Monitors process lifecycle

### Controller Module
- Coordinates all components
- Implements rebuild workflow
- Provides restart loop protection (>5 restarts in 10s triggers 5s throttle)
- Handles concurrent build requests

## Installation

```bash
# Clone the repository
git clone <repository-url>
cd hotreload

# Install dependencies
go mod download

# Build the tool
make build
```

## Usage

### Basic Command

```bash
hotreload --root <project-folder> --build "<build-command>" --exec "<run-command>"
```

### Example: Watch a Go Server

```bash
hotreload \
  --root ./myproject \
  --build "go build -o ./bin/server ./cmd/server" \
  --exec "./bin/server"
```

### Flags

- `--root` - Project root directory to watch (default: current directory)
- `--build` - Build command to execute (required)
- `--exec` - Command to run the server (required)

## Demo

The project includes a test HTTP server for demonstration:

```bash
# Run the demo (builds hotreload and watches testserver)
make demo

# Or use the alternative command
make run-demo
```

This will:
1. Build the hotreload tool
2. Start watching the `testserver` directory
3. Build and run the test server on http://localhost:8080
4. Automatically rebuild and restart when you edit `testserver/main.go`

Try it:
1. Run `make demo`
2. Visit http://localhost:8080 in your browser
3. Edit `testserver/main.go` (change the response message)
4. Save the file
5. Watch hotreload automatically rebuild and restart the server
6. Refresh your browser to see the changes

## How It Works

1. **Startup**: hotreload immediately builds and starts your server
2. **Watching**: Monitors all files in the project directory recursively
3. **Change Detection**: When files change, events are debounced (500ms window)
4. **Rebuild**: Executes the build command
5. **Restart**: If build succeeds, stops the old server and starts the new one
6. **Streaming**: Server logs appear in real-time with `[server]` prefix

## Real-World Problem Handling

### Multiple File Events
Editors like VSCode generate multiple events for a single save. The debouncer aggregates these into a single rebuild.

### Concurrent Changes
If files change during a build, the current build completes, then a new build starts with the latest state.

### Build Failures
Build errors are displayed prominently. The server is NOT started if the build fails.

### Process Termination
- First attempt: SIGTERM (graceful shutdown)
- After 5 seconds: SIGKILL (force kill)
- Process groups ensure child processes are also terminated

### Crash Loops
If the server restarts more than 5 times in 10 seconds, hotreload throttles for 5 seconds to prevent resource exhaustion.

### Dynamic Directories
New folders are automatically added to the watch list. Deleted folders are handled gracefully.

## Ignored Paths

The following are automatically ignored:

**Directories:**
- .git
- node_modules
- bin
- build
- dist
- tmp
- .idea
- .vscode
- vendor

**File Extensions:**
- .swp, .swx, .swo (Vim swap files)
- .tmp (temporary files)
- .log (log files)

**Hidden Files:**
- Files starting with `.` (except directories)

## Development

### Build

```bash
make build
```

### Run Tests

```bash
make test
```

### Format Code

```bash
make fmt
```

### Clean Build Artifacts

```bash
make clean
```

## Testing

Unit tests are provided for critical components:

```bash
# Run all tests
go test -v ./...

# Test specific module
go test -v ./internal/debounce
go test -v ./internal/filter
```

## Requirements

- Go 1.21 or later
- fsnotify library (automatically installed via go mod)

## Platform Support

- Linux
- macOS
- Windows

Process management is platform-aware and uses appropriate signals/methods for each OS.

## Logging

hotreload uses structured logging (log/slog) with the following levels:

- **Info**: Normal operations (build started, server started, file changes)
- **Warn**: Recoverable issues (graceful shutdown timeout, throttling)
- **Error**: Failures (build errors, process errors)
- **Debug**: Detailed information (individual file events)

## Performance

- **Restart time**: Typically 1-2 seconds (depends on build time)
- **Memory**: Minimal overhead (~10-20MB)
- **File descriptors**: Efficient watcher usage to avoid OS limits
- **Long-running**: Designed for hours of continuous operation

## Troubleshooting

### "Too many open files" error
Reduce the number of watched directories or increase your OS file descriptor limit:
```bash
# Linux/macOS
ulimit -n 4096
```

### Server doesn't stop
Check if your server handles signals properly. Implement graceful shutdown:
```go
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
go func() {
    <-sigChan
    // Cleanup and exit
    os.Exit(0)
}()
```

### Build command not found
Ensure the build command is in your PATH or use absolute paths.

## License

MIT

## Contributing

Contributions are welcome! Please ensure:
- Code follows Go best practices
- Tests are included for new features
- Documentation is updated
"# Trademarkia" 
