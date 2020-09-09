export GIT_COMMIT ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || echo unknown)
export VERSION    ?= ${GIT_COMMIT}
export LD_FLAGS = -X "main.commit=$(GIT_COMMIT)" -X "main.version=$(VERSION)" -X "main.date=$(shell date)"

export GOBIN = $(abspath .)/.tools/bin
export PATH := $(GOBIN):$(abspath .)/bin:$(PATH)
export CGO_ENABLED=0

export V = 0
export Q = $(if $(filter 1,$V),,@)
export M = $(shell printf "\033[34;1m▶\033[0m")

.PHONY: help
help:
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'

.PHONY: deps
deps:
	$(info $(M) fetching deps …)
	$Q go get -d -v ./...
	$Q go mod tidy

.PHONY: install-deptools
install-deptools: deps ## install dependent go tools
	$(info $(M) installing necessary tools …)
	$Q ./scripts/install_tools_check.sh

.PHONY: gen
gen:
	@go generate ./internal/templates/grpcwithgw

# Build servicebuilder binary
.PHONY: build
build: deps gen fmt vet
	go build -ldflags '$(LD_FLAGS)' -o bin/servicebuilder ./

# Run go fmt against code
.PHONY: fmt
fmt:
	go fmt ./...

# Run go vet against code
.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint: ## run golint
	$(info $(M) linting …)
	$Q ./.tools/bin/golangci-lint --verbose run ./... --timeout=5m

# Run tests
.PHONY: test
test: fmt vet
	go test ./... -coverprofile cover.out

.PHONY: clean
clean:
	@rm -rf bin dist
	@rm -rf test/tests.* test/coverage.*