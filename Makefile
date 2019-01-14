export PATH := ./bin:$(PATH)
export GO111MODULE := on

# Install all the build and lint dependencies
setup:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh
	go mod download
.PHONY: setup

# Build all files.
build:
	@echo "==> Building"
	@go build -o bin/giter
.PHONY: build

# Run all the linters
lint:
	@./bin/golangci-lint run
.PHONY: lint

# Release binaries to GitHub.
release: build
	@echo "==> Releasing"
	@goreleaser --rm-dist
	@echo "==> Complete"
.PHONY: release

# Clean.
clean:
	@rm -rf dist
.PHONY: clean
