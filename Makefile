NAME := k8sec

VERSION := $(patsubst "%",%,$(lastword $(shell grep "\tVersion" version.go)))
REVISION := $(shell git rev-parse --short HEAD)
BUILDTIME := $(shell date '+%Y/%m/%d %H:%M:%S %Z')
GOVERSION := $(subst go version ,,$(shell go version))

BINARYDIR := bin

LDFLAGS := -ldflags="-w -X \"main.GitCommit=$(REVISION)\" -X \"main.BuildTime=$(BUILDTIME)\" -X \"main.GoVersion=$(GOVERSION)\""

DISTDIR := dist

GITHUB_USERNAME := dtan4

.DEFAULT_GOAL := bin/$(NAME)

bin/$(NAME): deps
	go build $(LDFLAGS) -o bin/$(NAME)

build-all:
	go get github.com/mitchellh/gox
	gox -verbose \
		$(LDFLAGS) \
		-os="linux darwin windows " \
		-arch="amd64 386" \
		-output="$(DISTDIR)/{{.OS}}-{{.Arch}}/{{.Dir}}" .

.PHONY: clean
clean:
	rm -rf bin/*
	rm -rf dist/*
	rm -rf vendor/*

.PHONY: deps
deps: glide
	glide install

.PHONY: glide
glide:
ifeq ($(shell command -v glide 2> /dev/null),)
	curl https://glide.sh/get | sh
endif

.PHONY: install
install: deps
	go install $(LDFLAGS)

package-all:
	cd $(DISTDIR) \
	&& find * -type d | xargs -I {} tar czf $(BINARY)-$(VERSION)-{}.tar.gz {} \
	&& find * -type d | xargs -I {} zip -r $(BINARY)-$(VERSION)-{}.zip {}

.PHONY: test
test:
	go test -cover -v `glide novendor`

.PHONY: build-all package-all release-all
