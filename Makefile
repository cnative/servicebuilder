export GO111MODULE=on
export CGO_ENABLED=0
export GOBIN = $(abspath .)/.tools
export PATH := $(GOBIN):$(PATH)

export VERSION    ?= dev
export GIT_COMMIT ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || echo unknown)
export LD_FLAGS = -X "main.commit=$(GIT_COMMIT)" -X "main.version=$(VERSION)" -X "main.date=$(shell date)"

all: install-deptools clean test build

deps:
	go get -d -v ./...

install-deptools: deps
	@mkdir -p ${GOBIN}
	@go install github.com/go-bindata/go-bindata/...

gen:
	@go generate ./internal/templates/grpcwithgw

# Build servicebuilder binary
build: deps gen fmt vet
	go build -ldflags '$(LD_FLAGS)' -o bin/servicebuilder ./

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Run tests
test: fmt vet
	go test ./... -coverprofile cover.out

clean:
	@rm -rf bin dist
	@rm -rf test/tests.* test/coverage.*