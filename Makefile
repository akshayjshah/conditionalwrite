.DEFAULT_GOAL := all
SHELL := bash
.SHELLFLAGS := -euo pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --no-print-directory

.PHONY: help
help: ## Describe useful make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "%-30s %s\n", $$1, $$2}'

.PHONY: all
all: test lint ## Build, test, and lint (default)

.PHONY: test
test: ## Run tests
	go test -vet=off -race ./...

.PHONY: lint
lint: ## Lint Go
	test -z "$$(gofmt -s -l . | tee /dev/stderr)"
	go vet ./...
	go tool staticcheck ./...

.PHONY: lintfix
lintfix: ## Automatically fix some lint errors
	go run cmd/gofmt -s -w .

.PHONY: upgrade
upgrade: ## Upgrade Go dependencies
	go get -u -t ./...
	go mod tidy -v

