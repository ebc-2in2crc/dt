GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
NAME := dt
BUILDDIR=./build
BINDIR=$(BUILDDIR)/bin
PKGDIR=$(BUILDDIR)/pkg

VERSION := $(shell git describe --tags --abbrev=0)
LDFLAGS := -X 'main.version=$(VERSION)'
GOXOS := "darwin windows linux"
GOXARCH := "386 amd64"
GOXOUTPUT := "$(PKGDIR)/$(NAME)_{{.OS}}_{{.Arch}}/{{.Dir}}"

.PHONY: deps
## Install dependencies
deps:
	$(GOGET) golang.org/x/tools/cmd/goimports
	$(GOGET) github.com/golang/lint/golint
	$(GOGET) github.com/codegangsta/cli
	$(GOGET) github.com/mitchellh/go-homedir
	$(GOGET) github.com/Songmu/make2help/cmd/make2help

.PHONY: build
## Build binaries
build: deps
	go build -ldflags "$(LDFLAGS)" -o $(BINDIR)/$(NAME)

.PHONY: cross-build
## Cross build binaries
cross-build:
	rm -rf $(PKGDIR)
	gox -os=$(GOXOS) -arch=$(GOXARCH) -output=$(GOXOUTPUT)

.PHONY: test
## Run tests
test: deps
	$(GOTEST) -v ./...

.PHONY: lint
## Lint
lint: deps
	go vet ./...
	golint ./...

.PHONY: fmt
## Format source codes
fmt: deps
	find . -name "*.go" -not -path "./vendor/*" | xargs goimports -w

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -rf $(BUILDDIR)

.PHONY: help
## Show help
help:
	@make2help $(MAKEFILE_LIST)
