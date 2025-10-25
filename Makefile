GOPATH := $(shell go env GOPATH)
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
GOBIN := $(shell go env GOBIN)
RAINDROP_VERSION := "0.5.3"

GOOS ?= $(shell uname | tr '[:upper:]' '[:lower:]')
GOARCH ?=$(shell arch)

.PHONY: all build install

all: build install

.PHONY: mod-tidy
mod-tidy:
	go mod tidy

.PHONY: build OS ARCH
build: guard-RAINDROP_VERSION mod-tidy clean
	@echo "================================================="
	@echo "Building raindrop"
	@echo "================================================="

	@if [ ! -d "bin" ]; then \
		mkdir "bin"; \
	fi
	GOOS=${GOOS} GOARCH=${GOARCH} go build -o "bin/raindrop"
	sleep 2
	tar -czvf "raindrop_${RAINDROP_VERSION}_${GOOS}_${GOARCH}.tgz" bin; \

.PHONY: clean
clean:
	@echo "================================================="
	@echo "Cleaning raindrop"
	@echo "================================================="
	@if [ -f bin/raindrop ]; then \
		rm -f bin/raindrop; \
	fi; \

.PHONY: clean-all
clean-all: clean
	@echo "================================================="
	@echo "Cleaning tarballs"
	@echo "================================================="
	@rm -f *.tgz 2>/dev/null

.PHONY: install
install:
	@echo "================================================="
	@echo "Installing raindrop in ${GOPATH}/bin"
	@echo "================================================="

	GOOS=${GOOS} GOARCH=${GOARCH} go install

#
# General targets
#
guard-%:
	@if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi
