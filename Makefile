SHELL = /bin/sh

APP_NAME = funds
PACKAGES ?= ./...

GOBIN=bin
BINARY_PATH=$(GOBIN)/$(APP_NAME)

SERVER_SOURCE = cmd/api/api.go
SERVER_RUNNER = go run $(SERVER_SOURCE)
ifeq ($(DEBUG), true)
	SERVER_RUNNER = dlv debug $(SERVER_SOURCE) --
endif

## help: Display list of commands
.PHONY: help
help: Makefile
	@sed -n 's|^##||p' $< | column -t -s ':' | sed -e 's|^| |'

## $(APP_NAME)-api: Build app API with dependencies download
$(APP_NAME)-api: deps go build-api

## $(APP_NAME)-notifier: Build app Notifier with dependencies download
$(APP_NAME)-notifier: deps go build-notifier

## $(APP_NAME)-ui: Build app UI with dependencies download
$(APP_NAME)-ui: build-ui

.PHONY: go
go: format lint test bench

## name: Output name
.PHONY: name
name:
	@echo -n $(APP_NAME)

## dist: Output build output path
.PHONY: dist
dist:
	@echo -n $(BINARY_PATH)

## version: Output sha1 of last commit
.PHONY: version
version:
	@echo -n $(shell git rev-parse --short HEAD)

## author: Output author's name of last commit
.PHONY: author
author:
	@python -c 'import sys; import urllib; sys.stdout.write(urllib.quote_plus(sys.argv[1]))' "$(shell git log --pretty=format:'%an' -n 1)"

## deps: Download dependencies
.PHONY: deps
deps:
	go get github.com/kisielk/errcheck
	go get golang.org/x/lint/golint
	go get golang.org/x/tools/cmd/goimports

## format: Format code
.PHONY: format
format:
	goimports -w */*/*.go
	gofmt -s -w */*/*.go

## lint: Lint code
.PHONY: lint
lint:
	golint $(PACKAGES)
	errcheck -ignoretests $(PACKAGES)
	go vet $(PACKAGES)

## test: Test code with coverage
.PHONY: test
test:
	script/coverage

## bench: Benchmark code
.PHONY: bench
bench:
	go test $(PACKAGES) -bench . -benchmem -run Benchmark.*

## build-api: Build binary for api
.PHONY: build-api
build-api:
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo -o $(BINARY_PATH)-api $(SERVER_SOURCE)

## build-notifier: Build binary for notifier
.PHONY: build-notifier
build-notifier:
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo -o $(BINARY_PATH)-notifier cmd/alert/alert.go

## build-ui: Build bundle for ui
.PHONY: build-ui
build-ui:
	npm ci
	CI=true npm test
	npm run build

## start: Start app
.PHONY: start
start:
	$(SERVER_RUNNER)
