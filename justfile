# Goulder Dash - Boulder Dash clone in Go
# Available commands listed below

# Show this help message
default:
    @just --list

# Build the game binary
build:
    go build -o goulder-dash .

# Run the game directly
run:
    go run .

# Clean build artifacts
clean:
    rm -f goulder-dash
    go clean

# Format code using treefmt (gofumpt + gci)
fmt:
    treefmt --allow-missing-formatter

# Run golangci-lint
lint:
    golangci-lint run ./...

# Run golangci-lint
lint-fix:
    golangci-lint run ./... --fix

# Run go vet
vet:
    go vet ./...

# Run tests
test:
    go test ./...

# Run tests with verbose output
test-verbose:
    go test -v ./...

# Run tests with coverage
test-coverage:
    go test -cover ./...

# Check code quality (format + lint + vet + test)
check: fmt lint vet test

# Install development dependencies
install-deps:
    @echo "Installing development dependencies..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    go install mvdan.cc/gofumpt@latest
    go install github.com/daixiang0/gci@latest
    @echo "Installing treefmt..."
    @echo "Please install treefmt manually: https://github.com/numtide/treefmt"

# Tidy up go modules
mod-tidy:
    go mod tidy

# Update go modules
mod-update:
    go get -u ./...
    go mod tidy

# Build and run the game
dev: build run

# Continuous development - rebuild on file changes (requires entr)
watch:
    find . -name "*.go" | entr -r just run

# Create a release build with optimizations
release:
    go build -ldflags="-s -w" -o goulder-dash .

# Show project info
info:
    @echo "Project: Goulder Dash - Boulder Dash clone"
    @echo "Language: Go $(go version | cut -d' ' -f3)"
    @echo "Framework: gonutz/prototype"
    @echo ""
    @echo "Commands available:"
    @echo "  just build       - Build the game"
    @echo "  just run         - Run the game"
    @echo "  just dev         - Build and run"
    @echo "  just check       - Full quality check"
    @echo "  just fmt         - Format code"
    @echo "  just lint        - Run linter"
    @echo "  just test        - Run tests"

# Benchmark performance (if benchmarks exist)
bench:
    go test -bench=. ./...

# Profile the application (requires pprof-compatible code)
profile:
    go build -o goulder-dash .
    @echo "Run the game and visit http://localhost:6060/debug/pprof/"

# Check for security vulnerabilities
security:
    govulncheck ./...

# Generate documentation
docs:
    godoc -http=:6060
    @echo "Documentation server running at http://localhost:6060"