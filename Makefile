default: go docker

go: deps dev

dev: format lint tst bench build

docker: docker-deps docker-build

deps:
	go get -t ./...
	go get -u github.com/golang/lint/golint
	go get -u golang.org/x/tools/cmd/goimports

format:
	goimports -w **/*.go
	gofmt -s -w **/*.go

lint:
	golint ./...
	go vet ./...

tst:
	script/coverage

bench:
	go test ./... -bench . -benchmem -run Benchmark.*

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo -o bin/funds api.go
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo -o bin/notifier alert/alert.go

docker-deps:
	curl -s -o cacert.pem https://curl.haxx.se/ca/cacert.pem
	curl -s -o zoneinfo.zip https://raw.githubusercontent.com/golang/go/master/lib/time/zoneinfo.zip

docker-build:
	docker build -t ${DOCKER_USER}/funds-notifier -f alert/Dockerfile .
	docker build -t ${DOCKER_USER}/funds-front -f app/Dockerfile .
	docker build -t ${DOCKER_USER}/funds-api -f Dockerfile .

docker-push:
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASS}
	docker push ${DOCKER_USER}/funds-notifier
	docker push ${DOCKER_USER}/funds-front
	docker push ${DOCKER_USER}/funds-api
