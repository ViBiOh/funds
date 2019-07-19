SHELL = /bin/sh

ifneq ("$(wildcard .env)","")
	include .env
	export
endif

APP_NAME = funds
PACKAGES ?= ./...
GO_FILES ?= */*/*.go

GOBIN=bin
BINARY_PATH=$(GOBIN)/$(APP_NAME)

SERVER_SOURCE = cmd/api/api.go
SERVER_RUNNER = go run $(SERVER_SOURCE)
ifeq ($(DEBUG), true)
	SERVER_RUNNER = dlv debug $(SERVER_SOURCE) --
endif

NOTIFIER_SOURCE = cmd/alert/alert.go
NOTIFIER_RUNNER = go run $(NOTIFIER_SOURCE)
ifeq ($(DEBUG), true)
	NOTIFIER_RUNNER = dlv debug $(NOTIFIER_SOURCE) --
endif

.DEFAULT_GOAL := app

## help: Display list of commands
.PHONY: help
help: Makefile
	@sed -n 's|^##||p' $< | column -t -s ':' | sed -e 's|^| |'

## name: Output app name
.PHONY: name
name:
	@echo -n $(APP_NAME)

## version: Output last commit sha1
.PHONY: version
version:
	@echo -n $(shell git rev-parse --short HEAD)

## app: Build app with dependencies download
.PHONY: app
app: deps go

## go: Build app
.PHONY: go
go: format lint test bench build

## deps: Download dependencies
.PHONY: deps
deps:
	go get github.com/kisielk/errcheck
	go get golang.org/x/lint/golint
	go get golang.org/x/tools/cmd/goimports

## format: Format code
.PHONY: format
format:
	goimports -w $(GO_FILES)
	gofmt -s -w $(GO_FILES)

## lint: Lint code
.PHONY: lint
lint:
	golint $(PACKAGES)
	errcheck -ignoretests $(PACKAGES)
	go vet $(PACKAGES)

## test: Test with coverage
.PHONY: test
test:
	script/coverage

## bench: Benchmark code
.PHONY: bench
bench:
	go test $(PACKAGES) -bench . -benchmem -run Benchmark.*

## build: Build binary
.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo -o $(BINARY_PATH)-api $(SERVER_SOURCE)

## build-notifier: Build binary for notifier
.PHONY: build-notifier
build-notifier:
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo -o $(BINARY_PATH)-notifier $(NOTIFIER_SOURCE)

## start: Start app
.PHONY: start
start:
	$(SERVER_RUNNER)

## start: Start notifier
.PHONY: start-notifier
start-notifier:
	$(NOTIFIER_RUNNER) \
		-dbHost $(DATABASE_HOST) \
		-dbUser $(DATABASE_USER) \
		-dbPass $(DATABASE_PASS) \
		-dbName $(DATABASE_NAME) \
		-mailerURL https://mailer.vibioh.fr \
		-mailerUser $(MAILER_USER) \
		-mailerPass $(MAILER_PASS) \
		-recipients $(RECIPIENTS) \
		-score 20
