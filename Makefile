APP_NAME ?= funds
VERSION ?= $(shell git log --pretty=format:'%h' -n 1)
AUTHOR ?= $(shell git log --pretty=format:'%an' -n 1)

MAKEFLAGS += --silent
GOBIN=bin
BINARY_PATH=$(GOBIN)/$(APP_NAME)

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
go: format lint tst bench

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
	@echo -n $(VERSION)

## author: Output author's name of last commit
.PHONY: author
author:
	@python -c 'import sys; import urllib; sys.stdout.write(urllib.quote_plus(sys.argv[1]))' "$(AUTHOR)"

## deps: Download dependencies
.PHONY: deps
deps:
	go get github.com/golang/dep/cmd/dep
	go get github.com/kisielk/errcheck
	go get golang.org/x/lint/golint
	go get golang.org/x/tools/cmd/goimports
	dep ensure

## format: Format code
.PHONY: format
format:
	goimports -w */*/*.go
	gofmt -s -w */*/*.go

## lint: Lint code
.PHONY: lint
lint:
	golint `go list ./... | grep -v vendor`
	errcheck -ignoretests `go list ./... | grep -v vendor`
	go vet ./...

## tst: Test code with coverage
.PHONY: tst
tst:
	script/coverage

## bench: Benchmark code
.PHONY: bench
bench:
	go test ./... -bench . -benchmem -run Benchmark.*

## build-api: Build binary for api
.PHONY: build-api
build-api:
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo -o $(BINARY_PATH)-api cmd/api/api.go

## build-notifier: Build binary for notifier
.PHONY: build-notifier
build-notifier:
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo -o $(BINARY_PATH)-notifier cmd/alert/alert.go

## build-ui: Build bundle for ui
.PHONY: build-ui
build-ui:
	npm ci
	npm test
	npm run build

## start: Start app
.PHONY: start
start:
	go run cmd/api/api.go \
		-tls=false
