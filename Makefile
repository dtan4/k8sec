NAME := k8sec
VERSION := v0.6.0
COMMIT := $(shell git rev-parse HEAD)
DATE := $(shell date "+%Y-%m-%dT%H:%M:%S%z")

LDFLAGS  := -ldflags="-s -w -X \"github.com/dtan4/k8sec/version.version=$(VERSION)\" -X \"github.com/dtan4/k8sec/version.commit=$(COMMIT)\" -X \"github.com/dtan4/k8sec/version.date=$(DATE)\" -extldflags -static"

.DEFAULT_GOAL := build

.PHONY: build
build:
	go build $(LDFLAGS) -o bin/$(NAME)

.PHONY: clean
clean:
	rm -rf bin/*
	rm -rf dist/*
	rm -rf vendor/*

.PHONY: install
install:
	go install $(LDFLAGS)

.PHONY: test
test:
	go test -cover -race ./...
