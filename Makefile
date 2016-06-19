VRESION := 0.1.0
REVISION := $(shell git rev-parse --short HEAD)
BUILDTIME := $(shell date '+%Y/%m/%d %H:%M:%S %Z')
GOVERSION := $(subst go version ,,$(shell go version))

BINARYDIR := bin
BINARY := k8sec

LDFLAGS := -ldflags="-w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -X \"main.BuildTime=$(BUILDTIME)\" -X \"main.GoVersion=$(GOVERSION)\""

SOURCEDIR := .
SOURCES := $(shell find $(SOURCEDIR) -name '*.go' -type f)

GLIDE := glide
GLIDE_VERSION := 0.10.2

.DEFAULT_GOAL := $(BINARYDIR)/$(BINARY)

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

.PHONY: build
build:
	go build -ldflags="-w" -o $(BINARY_DIR)/$(BINARY)

.PHONY: clean
clean:
	rm -fr $(BINARYDIR)

.PHONY: deps
deps: $(BINARYDIR)/$(GLIDE)
	$(BINARYDIR)/$(GLIDE) install

.PHONY: install
install: $(BINARYDIR)/$(BINARY)
	go install
