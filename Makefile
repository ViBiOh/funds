SHELL := /bin/bash
DOCKER_VERSION ?= $(shell git log --pretty=format:'%h' -n 1)

default: go

go: deps api docker-build-api docker-push-api

api: format lint tst bench build-api

ui: node docker-build-ui docker-push-ui

notifier: deps format lint tst bench build-notifier

deps:
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/golang/lint/golint
	go get -u github.com/kisielk/errcheck
	go get -u golang.org/x/tools/cmd/goimports
	dep ensure

format:
	goimports -w **/*.go
	gofmt -s -w **/*.go

lint:
	golint `go list ./... | grep -v vendor`
	errcheck -ignoretests `go list ./... | grep -v vendor`
	go vet ./...

tst:
	script/coverage

bench:
	go test ./... -bench . -benchmem -run Benchmark.*

build-api:
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo -o bin/funds api.go

build-notifier:
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo -o bin/notifier alert/alert.go

node:
	npm run build

docker-deps:
	curl -s -o cacert.pem https://curl.haxx.se/ca/cacert.pem
	curl -s -o zoneinfo.zip https://raw.githubusercontent.com/golang/go/master/lib/time/zoneinfo.zip

docker-login:
	echo $(DOCKER_PASS) | docker login -u $(DOCKER_USER) --password-stdin

docker-promote: docker-promote-api docker-promote-notifier docker-promote-ui

docker-push: docker-push-api docker-push-notifier docker-push-ui

docker-build-api: docker-deps
	docker build -t $(DOCKER_USER)/funds-api:$(DOCKER_VERSION) -f Dockerfile .

docker-push-api: docker-login
	docker push $(DOCKER_USER)/funds-api

docker-promote-api:
	docker tag $(DOCKER_USER)/funds-api:$(DOCKER_VERSION) $(DOCKER_USER)/funds-api:latest

docker-build-ui: docker-deps
	docker build -t $(DOCKER_USER)/funds-ui:$(DOCKER_VERSION) -f app/Dockerfile .

docker-push-ui: docker-deps
	docker push $(DOCKER_USER)/funds-ui

docker-promote-ui:
	docker tag $(DOCKER_USER)/funds-ui:$(DOCKER_VERSION) $(DOCKER_USER)/funds-ui:latest

docker-build-notifier: docker-deps
	docker build -t $(DOCKER_USER)/funds-notifier:$(DOCKER_VERSION) -f alert/Dockerfile .

docker-push-notifier: docker-login
	docker push $(DOCKER_USER)/funds-notifier

docker-promote-notifier:
	docker tag $(DOCKER_USER)/funds-notifier:$(DOCKER_VERSION) $(DOCKER_USER)/funds-notifier:latest
