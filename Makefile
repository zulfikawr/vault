.PHONY: all build clean test install help run lint ui

BINARY_NAME=vault
LDFLAGS=-s -w

all: build

## build: Build the vault binary
build: ui
	CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) ./cmd/vault

## install: Install the binary to $GOPATH/bin
install:
	go install -ldflags "$(LDFLAGS)" ./cmd/vault

## run: Run the vault server
run:
	go run ./cmd/vault serve

## ui: Build the UI
ui:
	cd ui && bun install && bun x vite build

## clean: Remove binaries and build artifacts
clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -rf ui/dist

## test: Run tests
test:
	go test ./...

## lint: Run go fmt golangci-lint, go fmt, and eslint+prettier
lint:
	go fmt ./...
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run ./...
	cd ui && bun run lint

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
