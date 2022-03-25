MODULE   = $(shell env GO111MODULE=on $(GO) list -m)
DATE    ?= $(shell date +%FT%T%z)
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || \
			cat $(CURDIR)/.version 2> /dev/null || echo v0)
PKGS     = $(or $(PKG),$(shell env GO111MODULE=on $(GO) list ./...))
BIN      = $(CURDIR)/.bin

GOLANGCI_VERSION = v1.44.2

GO           = go
TIMEOUT_UNIT = 5m
TIMEOUT_E2E  = 20m
V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m🐱\033[0m")

export GO111MODULE=on

COMMANDS=$(patsubst cmd/%,%,$(wildcard cmd/*))
BINARIES=$(addprefix bin/,$(COMMANDS))

.PHONY: all
all: $(BINARIES) | $(BIN) ; $(info $(M) building executable...) @ ## Build program binary

$(BIN):
	@mkdir -p $@
$(BIN)/%: | $(BIN) ; $(info $(M) installing $(PACKAGE)...)
	$Q GOBIN=$(BIN) $(GO) install $(PACKAGE)

FORCE:

bin/%: cmd/% FORCE
	$Q $(GO) build -mod=vendor $(LDFLAGS) -v -o $@ ./$<

.PHONY: cross
cross: amd64 arm64  ## build cross platform binaries

.PHONY: amd64
amd64:
	GOOS=linux GOARCH=amd64 go build -mod=vendor $(LDFLAGS) ./cmd/...

.PHONY: arm64
arm64:
	GOOS=linux GOARCH=arm64 go build -mod=vendor $(LDFLAGS) ./cmd/...

KO = $(or ${KO_BIN},${KO_BIN},$(BIN)/ko)
$(BIN)/ko: PACKAGE=github.com/google/ko

.PHONY: apply
apply: | $(KO) ; $(info $(M) ko apply -R -f config/) @ ## Apply config to the current cluster
	$Q $(KO) apply -R -f config

.PHONY: resolve
resolve: | $(KO) ; $(info $(M) ko resolve -R -f config/) @ ## Resolve config to the current cluster
	$Q $(KO) resolve --push=false --oci-layout-path=$(BIN)/oci -R -f config

.PHONY: generated
generated: | vendor ; $(info $(M) update generated files) @ ## Update generated files
	$Q ./hack/update-codegen.sh

.PHONY: vendor
vendor: ; $(info $(M) update deps)
	$Q ./hack/update-deps.sh

## Tests
GOTESTSUM = $(BIN)/gotestsum
$(BIN)/gotestsum: PACKAGE=gotest.tools/gotestsum@v1.7.0

TEST_UNIT_TARGETS := test-unit-verbose test-unit-race
test-unit-verbose: ARGS=-v
test-unit-race:    ARGS=-race
$(TEST_UNIT_TARGETS): test-unit
.PHONY: $(TEST_UNIT_TARGETS) test-unit
test-unit: | $(GOTESTSUM) ; $(info $(M) run tests...) @ ## Run unit tests
	$Q $(GOTESTSUM) --format testname --packages="./..." -- -timeout $(TIMEOUT_UNIT) $(ARGS) $(COVER_OPTS)

TEST_E2E_TARGETS := test-e2e-short test-e2e-verbose test-e2e-race
test-e2e-short:   ARGS=-short
test-e2e-verbose: ARGS=-v
test-e2e-race:    ARGS=-race
$(TEST_E2E_TARGETS): test-e2e
.PHONY: $(TEST_E2E_TARGETS) test-e2e
test-e2e:  ## Run end-to-end tests
	$Q $(GOTESTSUM) --format testname --packages="./test/..." -- -timeout $(TIMEOUT_E2E) -tags e2e $(ARGS) $(COVER_OPTS)

.PHONY: check tests
check tests: test-unit test-e2e test-yamls

.PHONY: watch-test
watch-test: | $(GOTESTSUM) ; $(info $(M) watch and run tests) @ ## Watch and run tests
	$Q $(GOTESTSUM) --watch --format testname --packages="./..."

.PHONY: watch-resolve
watch-resolve: | $(KO) ; $(info $(M) watch and resolve config) @ ## Watch and build to the current cluster
	$Q $(KO) resolve -W --push=false --oci-layout-path=$(BIN)/oci -f config 1>/dev/null

.PHONY: watch-config
watch-config: | $(KO) ; $(info $(M) watch and apply config) @ ## Watch and apply to the current cluster
	$Q $(KO) apply -W -f config

# Linters and style
GOLANGCILINT = $(BIN)/golangci-lint
$(BIN)/golangci-lint: PACKAGE=github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_VERSION)

lint: | $(GOLANGCILINT) ; $(info $(M) running golangci-lint...) @ ## Run golangci-lint
	$Q $(GOLANGCILINT) run --modules-download-mode=vendor


GOFUMPT = $(BIN)/gofumpt
$(BIN)/gofumpt: PACKAGE=mvdan.cc/gofumpt@latest

.PHONY: fmt
fmt: | $(GOFUMPT) ; $(info $(M) running gofumpt...) @ ## Run gofumpt on all source files
	$Q $(GOFUMPT) -w -l cmd pkg test

# Misc

.PHONY: clean
clean: ; $(info $(M) cleaning...)	@ ## Cleanup everything
	@rm -rf $(BIN)
	@rm -rf bin
	@rm -rf test/tests.* coverage.*

.PHONY: help
help:
	@grep -hE '^[ a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-17s\033[0m %s\n", $$1, $$2}'

.PHONY: version
version:
	@echo $(VERSION)
