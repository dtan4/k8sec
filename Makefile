BINARY := k8sec
BINARY_DIR := bin

build:
	go build -o $(BINARY_DIR)/$(BINARY)

deps:
	go get github.com/Masterminds/glide
	glide install

install:
	go install

.PHONY: build deps install
