# Copyright (c) 2024 Six After, Inc.
#
# This source code is licensed under the Apache 2.0 License found in the
# LICENSE file in the root directory of this source tree.

SHELL := /bin/bash

.DEFAULT: ;: do nothing
GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_CLEAN=$(GO_CMD) clean
GO_TEST=$(GO_CMD) test
GO_GET=$(GO_CMD) get
GO_VET=$(GO_CMD) vet
GO_FMT=$(GO_CMD) fmt
GO_MOD=$(GO_CMD) mod
GO_LINT_CMD=golangci-lint run
GO_WORK=$(GO_CMD) work
GO_WORK_FILE := ./go.work
FUZZTIME ?= 20s

.PHONY: all
all: clean test

.PHONY: deps
deps: ## Get the dependencies and vendor
	@./scripts/go-deps.sh

.PHONY: test
test: ## Execute unit tests
	$(GO_TEST) -v ./...

.PHONY: fuzz
fuzz: ## Run each Go fuzz test individually (10s per test)
	@for fuzz in Fuzz_PRNG_Read Fuzz_NewReader; do \
		echo "===> Fuzzing $$fuzz"; \
		$(GO_TEST) -v -fuzz=$$fuzz -fuzztime=$(FUZZTIME) || exit $$?; \
	done

.PHONY: bench
bench: ## Execute benchmark tests for PRNG-CHACHA (raw bytes).
	@rm -f cpu.out
	@rm -f mem.out
	$(GO_TEST) -bench='^BenchmarkPRNG_' -run=^$$ -benchmem -memprofile=mem.out -cpuprofile=cpu.out .

.PHONY: bench-uuid
bench-uuid: ## Execute benchmark tests for using the PRNG-CHACHA to generate UUIDs using Google's uuid package.
	@rm -f cpu.out
	@rm -f mem.out
	$(GO_TEST) -bench='^BenchmarkUUID_' -run=^$$ -benchmem -memprofile=mem.out -cpuprofile=cpu.out .

.PHONY: clean
clean: ## Remove previous build
	$(GO_CLEAN) ./...

.PHONY: cover
cover: ## Generate global code coverage report
	@rm -f coverage.out
	$(GO_TEST) -v ./... -coverprofile coverage.out

.PHONY: analyze
analyze: ## Generate static analysis report
	$(GO_TEST) --json ./... -coverprofile coverage.out > coverage.json

.PHONY: fmt
fmt: ## Format the files
	$(GO_FMT) ./...

.PHONY: vet
vet: ## Vet the files
	$(GO_VET) -v ./...

.PHONY: lint
lint: ## Lint the files
	$(GO_LINT_CMD) --config .golangci.yaml --verbose ./...

.PHONY: tidy
tidy: ## Tidy vendored dependencies
	$(GO_MOD) tidy

.PHONY: vendor
vendor:
	@if [ -f $(GO_WORK_FILE) ]; then \
		$(GO_WORK) vendor; \
	else \
		$(GO_MOD) vendor; \
	fi

.PHONY: update
update: ## Update Go dependencies
	$(GO_GET) -u

.PHONY: vuln
vuln: ## Check for vulnerabilities
	govulncheck ./...

.PHONY: release-verify
release-verify: ## Verify the release
	rm -fr dist
	goreleaser --config .goreleaser.yaml release --snapshot

.PHONY: help
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' | sort

# %: - rule which match any task name;  @: - empty recipe = do nothing
%:
    @:
