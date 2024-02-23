.DEFAULT_GOAL := help

.PHONY: all
all: # Run all recipes.
all: lint test build

.PHONY: build
build: # Build a binary.
	@goreleaser --clean --skip publish --snapshot

.PHONY: help
help: # Show help information.
	@grep --extended-regexp "^[a-z-]+: #" "$(MAKEFILE_LIST)" | \
		awk 'BEGIN {FS = ": # "}; {printf "%-10s%s\n", $$1, $$2}'

.PHONY: lint
lint: # Lint the source code.
	@golangci-lint run

.PHONY: test
test: # Run tests.
	@go test -v ./...
