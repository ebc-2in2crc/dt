GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
NAME := dt
BUILDDIR=./build
BINDIR=$(BUILDDIR)/bin


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
	go build -o $(BINDIR)/$(NAME)

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
