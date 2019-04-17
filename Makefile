GOCMD := go
GOBUILD := $(GOCMD) build
GOINSTALL := $(GOCMD) install
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
NAME := dt
CURRENT := $(shell pwd)
CMDDIR := ./cmd/dt
BUILDDIR := ./build
BINDIR := $(BUILDDIR)/bin
PKGDIR := $(BUILDDIR)/pkg
DISTDIR := $(BUILDDIR)/dist

VERSION := $(shell git describe --tags --abbrev=0)
LDFLAGS := -X 'main.version=$(VERSION)'
GOXOS := "darwin windows linux"
GOXARCH := "386 amd64"
GOXOUTPUT := "$(PKGDIR)/$(NAME)_{{.OS}}_{{.Arch}}/{{.Dir}}"

.PHONY: deps
## Install dependencies
deps:
	$(GOGET) golang.org/x/tools/cmd/goimports
	$(GOGET) golang.org/x/lint/golint
	$(GOGET) github.com/codegangsta/cli
	$(GOGET) github.com/mitchellh/go-homedir
	$(GOGET) github.com/Songmu/make2help/cmd/make2help
	$(GOGET) github.com/tcnksm/ghr

.PHONY: build
## Build binaries
build: deps
	$(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BINDIR)/$(NAME) $(CMDDIR)

.PHONY: cross-build
## Cross build binaries
cross-build:
	rm -rf $(PKGDIR)
	gox -os=$(GOXOS) -arch=$(GOXARCH) -output=$(GOXOUTPUT) $(CMDDIR)

.PHONY: package
## Make package
package: cross-build
	rm -rf $(DISTDIR)
	mkdir $(DISTDIR)
	pushd $(PKGDIR) > /dev/null && \
		for P in `ls | xargs basename`; do zip -r $(CURRENT)/$(DISTDIR)/$$P.zip $$P; done && \
		popd > /dev/null

.PHONY: release
## Release package to Github
release: package
	ghr $(VERSION) $(DISTDIR)

.PHONY: install
## compile and install
install:
	$(GOINSTALL) -ldflags "$(LDFLAGS)" $(CMDDIR)

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
