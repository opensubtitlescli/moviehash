.DEFAULT_GOAL := help

.PHONY: all
all: # Run all recipes.
all: download lint test bin artifact

.PHONY: artifact
artifact: # Create an artifact.
	@tar \
		--create \
		--use-compress-program zstd \
		--directory .build \
		--file moviehash.tar.zst \
		moviehash

.PHONY: bin
bin: # Build a binary.
	@go build -o ./.build/moviehash ./cmd/moviehash/main.go

.PHONY: download
download: # Download modules.
	@go mod download

.PHONY: draft
draft: # Create a draft release.
	@v=$(shell ./.build/moviehash -v | tr -d "\n") && \
		gh release create "$$v" --draft --generate-notes && \
		gh release upload "$$v" moviehash.tar.zst

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
