BUILD_FILES = $(shell go list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}}\
{{end}}' ./...)
VERSION := $(shell git describe --tags $(shell git rev-list --tags --max-count=1))
DATE_FMT = +%Y-%m-%d
ifdef SOURCE_DATE_EPOCH
    BUILD_DATE ?= $(shell date -u -d "@$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>/dev/null || date -u -r "$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>/dev/null || date -u "$(DATE_FMT)")
else
    BUILD_DATE ?= $(shell date "$(DATE_FMT)")
endif

REVISION := $(shell git rev-parse --short HEAD)

ifndef CGO_CPPFLAGS
    export CGO_CPPFLAGS := $(CPPFLAGS)
endif
ifndef CGO_CFLAGS
    export CGO_CFLAGS := $(CFLAGS)
endif
ifndef CGO_LDFLAGS
    export CGO_LDFLAGS := $(LDFLAGS)
endif

GO_LDFLAGS := -X github.com/kmdkuk/git-push-notifier/version.Revision=$(REVISION) $(GO_LDFLAGS)
GO_LDFLAGS := -X github.com/kmdkuk/git-push-notifier/version.BuildDate=$(BUILD_DATE) $(GO_LDFLAGS)
DEV_LDFLAGS := $(GO_LDFLAGS)
GO_LDFLAGS := -X github.com/kmdkuk/git-push-notifier/version.Version=$(VERSION) $(GO_LDFLAGS)

# Test tools
BIN_DIR := $(shell pwd)/bin
GOLANGCI_LINT := $(BIN_DIR)/golangci-lint

bin/git-push-notifier: $(BUILD_FILES)
	go build -trimpath -ldflags "$(GO_LDFLAGS)" -o "$@" .

dev: $(BUILD_FILES)
	go build -trimpath -ldflags "$(DEV_LDFLAGS)" -o "bin/git-push-notifier-dev" .

install-go-tools:
	cat tools.go | awk -F'"' '/_/ {print $$2}' | xargs -tI {} go install {}
.PHONY: install-go-tools

test:
	./scripts/test_setup.sh
	go test -race -timeout 30m -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -html=coverage.txt -o coverage.html
.PHONY: test

lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run ./...
.PHONY: lint

profile:
	go tool pprof -http="localhost:8080" bin/git-push-notifier cpu.pprof
.PHONY: profile

$(GOLANGCI_LINT):
	GOBIN=$(BIN_DIR) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
