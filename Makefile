NAME := k8sec

VERSION := $(patsubst "%",%,$(lastword $(shell grep "\tVersion" version.go)))
REVISION := $(shell git rev-parse --short HEAD)
BUILDTIME := $(shell date '+%Y/%m/%d %H:%M:%S %Z')
GOVERSION := $(subst go version ,,$(shell go version))

BINARYDIR := bin

LDFLAGS := -ldflags="-w -X \"main.GitCommit=$(REVISION)\" -X \"main.BuildTime=$(BUILDTIME)\" -X \"main.GoVersion=$(GOVERSION)\""

GHR := ghr
GHR_VERSION := v0.4.0

DISTDIR := dist

GITHUB_USERNAME := dtan4

.DEFAULT_GOAL := bin/$(NAME)

$(BINARYDIR)/$(GHR):
ifeq ($(shell uname),Darwin)
	curl -fL https://github.com/tcnksm/ghr/releases/download/$(GHR_VERSION)/ghr_$(GHR_VERSION)_darwin_amd64.zip -o ghr.zip
	unzip ghr.zip
	if [ ! -d $(BINARYDIR) ]; then mkdir $(BINARYDIR); fi
	mv ./ghr $(BINARYDIR)/$(GHR)
	rm ./ghr.zip
else
	curl -fL https://github.com/tcnksm/ghr/releases/download/$(GHR_VERSION)/ghr-$(GHR_VERSION)-linux_amd64.zip -o ghr.zip
	unzip ghr.zip
	if [ ! -d $(BINARYDIR) ]; then mkdir $(BINARYDIR); fi
	mv ./ghr $(BINARYDIR)/$(GHR)
	rm ./ghr.zip
endif

bin/$(NAME): deps
	go build $(LDFLAGS) -o bin/$(NAME)

build-all:
	go get github.com/mitchellh/gox
	gox -verbose \
		$(LDFLAGS) \
		-os="linux darwin windows " \
		-arch="amd64 386" \
		-output="$(DISTDIR)/{{.OS}}-{{.Arch}}/{{.Dir}}" .

clean:
	rm -fr $(BINARYDIR)
	rm -fr $(DISTDIR)

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

release-all: build-all package-all $(BINARYDIR)/$(GHR)
	$(BINARYDIR)/$(GHR) -u $(GITHUB_USERNAME) --replace $(VERSION) $(DISTDIR)/

test:
	go test -v .

.PHONY: build-all clean package-all release-all test
