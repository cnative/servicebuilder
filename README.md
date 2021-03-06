# GO Micro Service Builder - project bootstrapping tool

[![CircleCI](https://circleci.com/gh/cnative/servicebuilder.svg?style=svg&circle-token=36641687cdfda3196b8569adb07ab4745117cdc4)](https://app.circleci.com/pipelines/github/cnative/servicebuilder)
[![Release](https://img.shields.io/github/release/cnative/servicebuilder.svg)](https://github.com/cnative/servicebuilder/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/cnative/servicebuilder)](https://goreportcard.com/report/github.com/cnative/servicebuilder)

Serivce Builder is a GO tool that generates scafolding / boilerplate code for building a cloud native service

A Micro Service generated by service builder

- enables rapid development of [gRPC](https://grpc.io/) based micro services
- exposes the gRPC services as REST / JSON via [grpc gateway](https://github.com/grpc-ecosystem/grpc-gateway) interface
- exposes metrics endpoint, which [Prometheus](https://prometheus.io/) could scrape from
- support tracing and metrics instrumentation using [OpenCensus](https://opencensus.io/)
- enables consistent logging
- exposes healthcheck endpoints for liveness and readiness probes
- defines state management interface and wraps them with metrics and tracing instrumentation automatically
- provides standard CLI
- postgres store for persistence
- build [Docker](https://www.docker.com/) container image
- enables [Kubernetes](https://kubernetes.io/) based deployment
  - [Helm chart](https://helm.sh/)
  - [K8S manifests](https://kustomize.io/) using [kustomize](https://github.com/kubernetes-sigs/kustomize)

## Getting Started

### Pre-Req

- [Go 1.15 +](https://golang.org/dl/)

### Install

#### `homebrew` MacOS

```bash
brew tap cnative/tap
brew install servicebuilder
```

#### `go get`

```bash
GOBIN=$HOME/bin go get -u github.com/cnative/servicebuilder
```

if your `$HOME/bin` is not already in `PATH`

```bash
export $HOME/bin:$PATH
```

`Note:` Service builder downloaded via `go get` will download the latest version of service builder available at the time of execution. If a fixed version is needed download a [pre-built binary](https://github.com/cnative/servicebuider/releases)

#### Pre-built binary

The easiest way to get `servicebuilder` is to use one of the [pre-built release](https://github.com/cnative/servicebuilder/releases) binaries which are available for OSX, Linux, and Windows on the release page

### Building from source

To check out this repository:

- Create your own [fork](https://help.github.com/articles/fork-a-repo/) of `github.com/cnative/servicebuilder`
- Clone it at some location on your host

```bash
git clone git@github.com:{YOUR_GITHUB_USERNAME}/servicebuilder
cd servicebuilder
git remote add upstream git@github.com:cnative/servicebuilder.git
git remote set-url --push upstream no_push


```
=======
make install-deptools clean build
```
