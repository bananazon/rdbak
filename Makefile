GOPATH := $(shell go env GOPATH)
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
GOBIN := $(shell go env GOBIN)
rdbak_VERSION := "0.3.1"

GOOS ?= $(shell uname | tr '[:upper:]' '[:lower:]')
GOARCH ?=$(shell arch)

.PHONY: all build install

all: build install

.PHONY: mod-tidy
mod-tidy:
	go mod tidy

.PHONY: build OS ARCH
build: guard-rdbak_VERSION mod-tidy clean
	@echo "================================================="
	@echo "Building rdbak"
	@echo "================================================="

	@if [ ! -d "bin" ]; then \
		mkdir "bin"; \
	fi
	GOOS=${GOOS} GOARCH=${GOARCH} go build -o "bin/rdbak"
	sleep 2
	tar -czvf "rdbak_${rdbak_VERSION}_${GOOS}_${GOARCH}.tgz" bin; \

.PHONY: clean
clean:
	@echo "================================================="
	@echo "Cleaning rdbak"
	@echo "================================================="
	@if [ -f bin/rdbak ]; then \
		rm -f bin/rdbak; \
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
	@echo "Installing rdbak in ${GOPATH}/bin"
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
