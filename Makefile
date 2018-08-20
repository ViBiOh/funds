APP_NAME := funds
VERSION ?= $(shell git log --pretty=format:'%h' -n 1)
AUTHOR ?= $(shell git log --pretty=format:'%an' -n 1)

$(APP_NAME)-api: deps go build-api

$(APP_NAME)-notifier: deps go build-notifier

$(APP_NAME)-ui: build-ui

go: format lint tst bench

name:
	@echo -n $(APP_NAME)

version:
	@echo -n $(VERSION)

author:
	@python -c 'import sys; import urllib; sys.stdout.write(urllib.quote_plus(sys.argv[1]))' "$(AUTHOR)"

deps:
	go get github.com/golang/dep/cmd/dep
	go get github.com/golang/lint/golint
	go get github.com/kisielk/errcheck
	go get golang.org/x/tools/cmd/goimports
	dep ensure

format:
	goimports -w */*/*.go
	gofmt -s -w */*/*.go

lint:
	golint `go list ./... | grep -v vendor`
	errcheck -ignoretests `go list ./... | grep -v vendor`
	go vet ./...

tst:
	script/coverage

bench:
	go test ./... -bench . -benchmem -run Benchmark.*

build-api:
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo -o bin/$(APP_NAME)-api cmd/api/api.go

build-notifier:
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo -o bin/$(APP_NAME)-notifier cmd/alert/alert.go

build-ui:
	npm ci
	npm run build

start:
	go run cmd/api/api.go \
		-tls=false

.PHONY: $(APP_NAME)-api $(APP_NAME)-notifier $(APP_NAME)-ui go name version author deps format lint tst bench build-api build-notifier build-ui start
