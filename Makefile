.PHONY: all build fmt vet lint test tidy check clean

# Default: run all checks
all: check build

# Build the binary
build:
	go build -o gitlab-mcp-server .

# Format code
fmt:
	gofmt -s -w .

# Run go vet
vet:
	go vet ./...

# Run golangci-lint
lint:
	golangci-lint run ./...

# Run tests
test:
	go test -v -race ./...

# Tidy modules
tidy:
	go mod tidy
	@git diff --exit-code go.mod go.sum || (echo "go.mod/go.sum not tidy" && exit 1)

# Run all checks (used by pre-commit and CI)
check: fmt vet lint tidy
	@echo "All checks passed"

# Clean build artifacts
clean:
	rm -f gitlab-mcp-server
