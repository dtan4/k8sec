NAME := k8sec
VERSION := v0.6.0
REVISION := $(shell git rev-parse --short HEAD)

SRCS     := $(shell find . -type f -name '*.go')
LDFLAGS  := -ldflags="-s -w -X \"github.com/dtan4/k8sec/version.Version=$(VERSION)\" -X \"github.com/dtan4/k8sec/version.Revision=$(REVISION)\" -extldflags -static"
NOVENDOR := $(shell go list ./... | grep -v vendor)

DIST_DIRS := find * -type d -exec

DOCKER_REPOSITORY := quay.io
DOCKER_IMAGE_NAME := $(DOCKER_REPOSITORY)/dtan4/k8sec
DOCKER_IMAGE_TAG  ?= latest
DOCKER_IMAGE      := $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

export GO111MODULE=on

.DEFAULT_GOAL := bin/$(NAME)

bin/$(NAME): $(SRCS)
	go build $(LDFLAGS) -o bin/$(NAME)

.PHONY: ci-docker-release
ci-docker-release: docker-build
	@docker login -u="$(DOCKER_QUAY_USERNAME)" -p="$(DOCKER_QUAY_PASSWORD)" $(DOCKER_REPOSITORY)
	docker push $(DOCKER_IMAGE)

.PHONY: ci-test
ci-test:
	echo "" > coverage.txt
	set -e; \
	for d in $(NOVENDOR); do \
		go test -coverprofile=profile.out -covermode=atomic -v $$d; \
		if [ -f profile.out ]; then \
			cat profile.out >> coverage.txt; \
			rm profile.out; \
		fi; \
	done

.PHONY: clean
clean:
	rm -rf bin/*
	rm -rf dist/*
	rm -rf vendor/*

.PHONY: cross-build
cross-build:
	for os in darwin linux windows; do \
		for arch in amd64 386; do \
			GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o dist/$$os-$$arch/$(NAME); \
		done; \
	done
.PHONY: dist
dist:
	cd dist && \
	$(DIST_DIRS) cp ../LICENSE {} \; && \
	$(DIST_DIRS) cp ../README.md {} \; && \
	$(DIST_DIRS) tar -zcf $(NAME)-$(VERSION)-{}.tar.gz {} \; && \
	$(DIST_DIRS) zip -r $(NAME)-$(VERSION)-{}.zip {} \; && \
	cd ..

.PHONY: docker-build
docker-build:
	docker build -t $(DOCKER_IMAGE) .

.PHONY: fast
fast:
	go build $(LDFLAGS) -o bin/$(NAME)

.PHONY: install
install:
	go install $(LDFLAGS)

.PHONY: release
release:
	git tag $(VERSION)
	git push origin $(VERSION)

.PHONY: test
test:
	go test -cover -v $(NOVENDOR)
