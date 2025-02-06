VERSION ?= $(shell git describe --tags --always --dirty)
BRANCH  ?= $(shell git rev-parse --abbrev-ref HEAD)
COMMIT  ?= $(shell git rev-parse --short HEAD)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d %H:%M:%S')

BINARY_NAME=solana-program-scanner
LDFLAGS := -s -w \
           -X 'main.Version=$(VERSION)' \
           -X 'main.GitBranch=$(BRANCH)' \
           -X 'main.GitCommit=$(COMMIT)' \
           -X 'main.BuildTime=$(BUILD_TIME)'

.PHONY: test clean release

test:
	go build -ldflags "$(LDFLAGS)"
	./$(BINARY_NAME) run > scanner.log 2>&1

clean:
	rm -rf $(BINARY_NAME) scanner.log blocks.json

release:
	# Build for Linux
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-linux-amd64
	tar czf $(BINARY_NAME)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64

	# Build for macOS
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-darwin-amd64
	tar czf $(BINARY_NAME)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64

	# Clean up binaries
	rm -f $(BINARY_NAME)-linux-amd64 $(BINARY_NAME)-darwin-amd64

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)
