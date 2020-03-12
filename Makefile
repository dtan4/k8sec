NAME := k8sec
VERSION := v0.6.0
COMMIT := $(shell git rev-parse HEAD)
DATE := $(shell date "+%Y-%m-%dT%H:%M:%S%z")

SRCS     := $(shell find . -type f -name '*.go')
LDFLAGS  := -ldflags="-s -w -X \"github.com/dtan4/k8sec/version.version=$(VERSION)\" -X \"github.com/dtan4/k8sec/version.commit=$(COMMIT)\" -X \"github.com/dtan4/k8sec/version.date=$(DATE)\" -extldflags -static"
NOVENDOR := $(shell go list ./... | grep -v vendor)

export GO111MODULE=on

.DEFAULT_GOAL := bin/$(NAME)

bin/$(NAME): $(SRCS)
	go build $(LDFLAGS) -o bin/$(NAME)

.PHONY: ci-test
ci-test:
	go test -coverpkg=./... -coverprofile=coverage.txt -v ./...

.PHONY: clean
clean:
	rm -rf bin/*
	rm -rf dist/*
	rm -rf vendor/*

.PHONY: fast
fast:
	go build $(LDFLAGS) -o bin/$(NAME)

.PHONY: install
install:
	go install $(LDFLAGS)

.PHONY: test
test:
	go test -cover -v $(NOVENDOR)
