# Load variables from .env.local if present
ifneq (,$(wildcard .env.local))
include .env.local
export $(shell sed -n 's/^\([A-Za-z_][A-Za-z0-9_]*\)=.*/\1/p' .env.local)
endif

# Default target
.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help.
	@awk 'BEGIN {FS = ":.*?## "}; /^[a-zA-Z0-9_\-]+:.*?## / {i++; T[i] = $$1; D[i] = $$2; if (length($$1) > max_len) max_len = length($$1)} END {max_len += 2; for (j = 1; j <= i; j++) printf "\033[36m%-" max_len "s\033[0m %s\n", T[j], D[j]}' $(MAKEFILE_LIST)

.PHONY: build-debug
build-debug: ## Build with debug gcflags.
	@go build -gcflags "all=-N -l" -o gif2vid ./cmd/gif2vid

.PHONY: local-debug
local-debug: build-debug ## Build and run with debugger.
	@dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./gif2vid

.PHONY: build
build: ## Build to gif2vid.
	@go build -o tmp/gif2vid ./cmd/gif2vid

.PHONY: local
local: build ## Build and run binary locally.
	@./tmp/gif2vid

.PHONY: install
install: ## Run go install.
	@go install ./cmd/gif2vid

.PHONY: deps
deps: ## Run go mod tidy and vendor.
	@go mod tidy
	@go mod vendor

.PHONY: fmt
fmt: ## Run go fmt.
	@go fmt ./...

.PHONY: vet
vet: ## Run go vet.
	@go vet ./...

.PHONY: test
test: ## Run go test.
	@go test ./...
