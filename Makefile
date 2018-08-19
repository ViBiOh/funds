APP_NAME := funds
VERSION ?= $(shell git log --pretty=format:'%h' -n 1)
AUTHOR ?= $(shell git log --pretty=format:'%an' -n 1)

docker:
	docker build -t vibioh/$(APP_NAME)-api:$(VERSION) .

notifier:
	docker build -t vibioh/$(APP_NAME)-notifier:$(VERSION) -f Dockerfile_notifier .

ui:
	docker build -t vibioh/$(APP_NAME)-ui:$(VERSION) -f Dockerfile_ui .

$(APP_NAME)-api: deps go

$(APP_NAME)-notifier: deps build-notifier

go: format lint tst bench build-api

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

start:
	go run cmd/api/api.go \
		-tls=false

.PHONY: docker notifier ui $(APP_NAME)-api $(APP_NAME)-notifier go name version author deps format lint tst bench build-api build-notifier start
