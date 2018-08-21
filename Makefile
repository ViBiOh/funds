MAKEFLAGS += --silent
GOBIN=bin
BINARY_PATH=$(GOBIN)/$(APP_NAME)
VERSION ?= $(shell git log --pretty=format:'%h' -n 1)
AUTHOR ?= $(shell git log --pretty=format:'%an' -n 1)

APP_NAME ?= funds

help: Makefile
	@sed -n 's|^##||p' $< | column -t -s ':' | sed -e 's|^| |'

## $(APP_NAME)-api: Build app API with dependencies download
$(APP_NAME)-api: deps go build-api

## $(APP_NAME)-notifier: Build app Notifier with dependencies download
$(APP_NAME)-notifier: deps go build-notifier

## $(APP_NAME)-ui: Build app UI with dependencies download
$(APP_NAME)-ui: build-ui

go: format lint tst bench

## name: Output name of app
name:
	@echo -n $(APP_NAME)

## dist: Output build output path
dist:
	@echo -n $(BINARY_PATH)

## version: Output sha1 of last commit
version:
	@echo -n $(VERSION)

## author: Output author's name of last commit
author:
	@python -c 'import sys; import urllib; sys.stdout.write(urllib.quote_plus(sys.argv[1]))' "$(AUTHOR)"

## deps: Download dependencies
deps:
	go get github.com/golang/dep/cmd/dep
	go get github.com/golang/lint/golint
	go get github.com/kisielk/errcheck
	go get golang.org/x/tools/cmd/goimports
	dep ensure

## format: Format code of app
format:
	goimports -w */*/*.go
	gofmt -s -w */*/*.go

## lint: Lint code of app
lint:
	golint `go list ./... | grep -v vendor`
	errcheck -ignoretests `go list ./... | grep -v vendor`
	go vet ./...

## tst: Test code of app with coverage
tst:
	script/coverage

## bench: Benchmark code of app
bench:
	go test ./... -bench . -benchmem -run Benchmark.*

## build-api: Build binary of app
build-api:
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo -o $(BINARY_PATH)-api cmd/api/api.go

## build-notifier: Build binary of app
build-notifier:
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo -o $(BINARY_PATH)-notifier cmd/alert/alert.go

## build-ui: Build bundle of app
build-ui:
	npm ci
	npm run build

## start: Start app
start:
	go run cmd/api/api.go \
		-tls=false

.PHONY: help $(APP_NAME)-api $(APP_NAME)-notifier $(APP_NAME)-ui go name dist version author deps format lint tst bench build-api build-notifier build-ui start
