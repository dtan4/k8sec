VERSION := $(patsubst "%",%,$(lastword $(shell grep "\tVersion" version.go)))
REVISION := $(shell git rev-parse --short HEAD)
BUILDTIME := $(shell date '+%Y/%m/%d %H:%M:%S %Z')
GOVERSION := $(subst go version ,,$(shell go version))

BINARYDIR := bin
BINARY := k8sec

LDFLAGS := -ldflags="-w -X \"main.GitCommit=$(REVISION)\" -X \"main.BuildTime=$(BUILDTIME)\" -X \"main.GoVersion=$(GOVERSION)\""

SOURCEDIR := .
SOURCES := $(shell find $(SOURCEDIR) -name '*.go' -type f)

GHR := ghr
GHR_VERSION := v0.4.0

GLIDE := glide
GLIDE_VERSION := 0.10.2

DISTDIR := dist

GITHUB_USERNAME := dtan4

.DEFAULT_GOAL := $(BINARYDIR)/$(BINARY)

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

$(BINARYDIR)/$(GLIDE):
ifeq ($(shell uname),Darwin)
	curl -fL https://github.com/Masterminds/glide/releases/download/$(GLIDE_VERSION)/glide-$(GLIDE_VERSION)-darwin-amd64.zip -o glide.zip
	unzip glide.zip
	if [ ! -d $(BINARYDIR) ]; then mkdir $(BINARYDIR); fi
	mv ./darwin-amd64/glide $(BINARYDIR)/$(GLIDE)
	rm -fr ./darwin-amd64
	rm ./glide.zip
else
	curl -fL https://github.com/Masterminds/glide/releases/download/$(GLIDE_VERSION)/glide-$(GLIDE_VERSION)-linux-amd64.zip -o glide.zip
	unzip glide.zip
	if [ ! -d $(BINARYDIR) ]; then mkdir $(BINARYDIR); fi
	mv ./linux-amd64/glide $(BINARYDIR)/$(GLIDE)
	rm -fr ./linux-amd64
	rm ./glide.zip
endif

$(BINARYDIR)/$(BINARY): $(SOURCES)
	go build $(LDFLAGS) -o $(BINARYDIR)/$(BINARY)

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

deps: $(BINARYDIR)/$(GLIDE)
	$(BINARYDIR)/$(GLIDE) install

install: $(BINARYDIR)/$(BINARY)
	go install

package-all:
	cd $(DISTDIR) \
	&& find * -type d | xargs -I {} tar czf $(BINARY)-$(VERSION)-{}.tar.gz {} \
	&& find * -type d | xargs -I {} zip -r $(BINARY)-$(VERSION)-{}.zip {}

release-all: build-all package-all $(BINARYDIR)/$(GHR)
	$(BINARYDIR)/$(GHR) -u $(GITHUB_USERNAME) --delete --replace $(VERSION) $(DISTDIR)/

test:
	go test -v .

.PHONY: build-all clean deps install package-all release-all test
