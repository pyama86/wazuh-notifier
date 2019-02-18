INFO_COLOR=\033[1;34m
RESET=\033[0m
BOLD=\033[1m
GOVERSION ?= $(shell go version | awk '{print $$3}')
ifeq ("$(shell uname)","Darwin")
GO ?= GO111MODULE=on go
else
GO ?= GO111MODULE=on /usr/local/go/bin/go
endif

TEST ?= $(shell go list ./... | grep -v -e vendor -e keys -e tmp)
VERSION = $(shell cat version)

test: ## Run test
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Testing$(RESET)"
	$(GO) test -v $(TEST) -timeout=30s -parallel=4
	$(GO) test -race $(TEST)

build: ## Build server
	$(GO) build -ldflags "-s -w -X main.goversion=$(GOVERSION)" -o wazuh-notifier
