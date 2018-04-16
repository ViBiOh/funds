VERSION ?= $(shell git log --pretty=format:'%h' -n 1)
APP_NAME := funds

default: api

api: deps go docker-api

go: format lint tst bench build-api

notifier: deps build-notifier docker-notifier

version:
	@echo -n $(VERSION)

deps:
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/golang/lint/golint
	go get -u github.com/kisielk/errcheck
	go get -u golang.org/x/tools/cmd/goimports
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
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo -o bin/$(APP_NAME) cmd/api/api.go

build-notifier:
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo -o bin/notifier cmd/alert/alert.go

docker-deps:
	curl -s -o cacert.pem https://curl.haxx.se/ca/cacert.pem
	curl -s -o zoneinfo.zip https://raw.githubusercontent.com/golang/go/master/lib/time/zoneinfo.zip

docker-login:
	echo $(DOCKER_PASS) | docker login -u $(DOCKER_USER) --password-stdin

docker-pull: docker-pull-api docker-pull-ui docker-pull-notifier

docker-promote: docker-pull docker-promote-api docker-promote-ui docker-promote-notifier

docker-delete: docker-delete-api docker-delete-ui docker-delete-notifier

docker-push: docker-push-api docker-push-ui docker-push-notifier

docker-api: docker-build-api docker-push-api

docker-ui: docker-build-ui docker-push-ui

docker-notifier: docker-build-notifier docker-push-notifier

docker-build-api: docker-deps
	docker build -t $(DOCKER_USER)/$(APP_NAME)-api:$(VERSION) -f Dockerfile .

docker-push-api: docker-login
	docker push $(DOCKER_USER)/$(APP_NAME)-api:$(VERSION)

docker-pull-api:
	docker pull $(DOCKER_USER)/$(APP_NAME)-api:$(VERSION)

docker-promote-api:
	docker tag $(DOCKER_USER)/$(APP_NAME)-api:$(VERSION) $(DOCKER_USER)/$(APP_NAME)-api:latest

docker-delete-api:
	curl -X DELETE -u "$(DOCKER_USER):$(DOCKER_CLOUD_TOKEN)" "https://cloud.docker.com/v2/repositories/$(DOCKER_USER)/$(APP_NAME)-api/tags/$(VERSION)/"

docker-build-ui: docker-deps
	docker build -t $(DOCKER_USER)/$(APP_NAME)-ui:$(VERSION) -f ui/Dockerfile ./ui/

docker-push-ui: docker-login
	docker push $(DOCKER_USER)/$(APP_NAME)-ui:$(VERSION)

docker-pull-ui:
	docker pull $(DOCKER_USER)/$(APP_NAME)-ui:$(VERSION)

docker-promote-ui:
	docker tag $(DOCKER_USER)/$(APP_NAME)-ui:$(VERSION) $(DOCKER_USER)/$(APP_NAME)-ui:latest

docker-delete-ui:
	curl -X DELETE -u "$(DOCKER_USER):$(DOCKER_CLOUD_TOKEN)" "https://cloud.docker.com/v2/repositories/$(DOCKER_USER)/$(APP_NAME)-ui/tags/$(VERSION)/"

docker-build-notifier: docker-deps
	docker build -t $(DOCKER_USER)/$(APP_NAME)-notifier:$(VERSION) -f cmd/alert/Dockerfile .

docker-push-notifier: docker-login
	docker push $(DOCKER_USER)/$(APP_NAME)-notifier:$(VERSION)

docker-pull-notifier:
	docker pull $(DOCKER_USER)/$(APP_NAME)-notifier:$(VERSION)

docker-promote-notifier:
	docker tag $(DOCKER_USER)/$(APP_NAME)-notifier:$(VERSION) $(DOCKER_USER)/$(APP_NAME)-notifier:latest

docker-delete-notifier:
	curl -X DELETE -u "$(DOCKER_USER):$(DOCKER_CLOUD_TOKEN)" "https://cloud.docker.com/v2/repositories/$(DOCKER_USER)/$(APP_NAME)-notifier/tags/$(VERSION)/"

.PHONY: api go notifier version deps format lint tst bench build-api build-notifier docker-deps docker-login docker-pull docker-promote docker-push docker-api docker-ui docker-notifier docker-build-api docker-push-api docker-pull-api docker-promote-api docker-build-ui docker-push-ui docker-pull-ui docker-promote-ui docker-build-notifier docker-push-notifier docker-pull-notifier docker-promote-notifier
