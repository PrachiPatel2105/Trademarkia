.PHONY: build run-demo demo test clean

# Build the hotreload CLI tool
build:
	go build -o bin/hotreload.exe ./cmd/hotreload

# Build the test server
build-testserver:
	go build -o bin/testserver.exe ./testserver

# Run demo: build hotreload and use it to watch the test server
demo: build
	./bin/hotreload.exe --root ./testserver --build "go build -o ./bin/testserver.exe ./testserver" --exec "./bin/testserver.exe"

# Alternative demo command
run-demo: demo

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	golangci-lint run
