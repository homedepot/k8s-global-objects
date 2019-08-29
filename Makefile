export GO111MODULE=on

MAKEFLAGS += --silent
.DEFAULT_GOAL := help

# Project name is the same as the binary name in .goreleaser.yml
PROJECTNAME := k8s-global-objects
PROJECTORG := homedepot
GOPATH := $(GOPATH)

GORELEASER_VERSION := 0.99.0
TAG := $(shell cat version/version.go | grep "Version" | head -1 | sed 's/\"//g' | cut -d' ' -f3 )

## build: Build local binaries and docker image. Requires `go` to be installed.
build: | format test
	@echo "=> Building with goreleaser ..."
	git tag v$(TAG) --force
	goreleaser --skip-publish --skip-validate
	make clean
.PHONY: build

format:
	@echo "=> Running Go Format ..."
	go list ./... | grep -v /vendor/ | xargs go vet -composites=false
	find ./ -name "*.go" | grep -v /vendor/ | xargs goimports -w=true
.PHONY: format

test:
	@echo "=> Running Go Test via Overalls ..."
	mkdir _test
	go get golang.org/x/tools/cmd/cover
	go get github.com/go-playground/overalls
	overalls -covermode=atomic -project=github.com/"$(PROJECTORG)"/"$(PROJECTNAME)" -- -race -v
	mv overalls.coverprofile _test/$(PROJECTNAME).cover
	go tool cover -func=_test/$(PROJECTNAME).cover
	make clean
.PHONY: test

## install-goreleaser-linux: Install goreleaser on your system for Linux systems.
install-goreleaser-linux:
	wget -nv -P /tmp/ https://github.com/goreleaser/goreleaser/releases/download/v$(GORELEASER_VERSION)/goreleaser_Linux_x86_64.tar.gz
	tar -C ~/bin -xzf /tmp/goreleaser_Linux_x86_64.tar.gz goreleaser
	rm -r /tmp/goreleaser_Linux_x86_64.tar.gz
.PHONY: install-goreleaser-linux

## install-goreleaser-darwin: Install goreleaser on your system for macOS (Darwin).
install-goreleaser-darwin:
	wget -nv -P /tmp/ https://github.com/goreleaser/goreleaser/releases/download/v$(GORELEASER_VERSION)/goreleaser_Darwin_x86_64.tar.gz
	tar -C /usr/local/bin -xzf /tmp/goreleaser_Darwin_x86_64.tar.gz goreleaser
	rm -r /tmp/goreleaser_Darwin_x86_64.tar.gz
.PHONY: install-goreleaser-darwin

## github-release: Publish a release to github.
github-release: | test
	@echo "=> Running Publish Release to Github ..."
	git tag v$(TAG) --force
	git push origin v$(TAG) --force
	goreleaser
.PHONY: github-release

## clean: Clean directory.
clean:
	@echo "=> Cleaning directories ..."
	rm -rf _dist/
	rm -rf _test/
	rm -rf runner/*coverprofile
.PHONY: clean

all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
.PHONY: help
