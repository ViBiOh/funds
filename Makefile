default: go docker

go: deps dev

dev: format lint tst bench build

docker: docker-deps docker-build

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

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo -o bin/funds api.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-s -w" -installsuffix nocgo -o bin/funds-arm api.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -installsuffix nocgo -o bin/funds-arm64 api.go
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo -o bin/notifier alert/alert.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-s -w" -installsuffix nocgo -o bin/notifier-arm alert/alert.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -installsuffix nocgo -o bin/notifier-arm64 alert/alert.go

docker-deps:
	curl -s -o cacert.pem https://curl.haxx.se/ca/cacert.pem
	curl -s -o zoneinfo.zip https://raw.githubusercontent.com/golang/go/master/lib/time/zoneinfo.zip

docker-build:
	docker build -t ${DOCKER_USER}/funds-notifier -f alert/Dockerfile .
	docker build -t ${DOCKER_USER}/funds-notifier:arm -f alert/Dockerfile_arm .
	docker build -t ${DOCKER_USER}/funds-notifier:arm64 -f alert/Dockerfile_arm64 .
	docker build -t ${DOCKER_USER}/funds-front -f app/Dockerfile .
	docker build -t ${DOCKER_USER}/funds-front:arm -f app/Dockerfile_arm .
	docker build -t ${DOCKER_USER}/funds-front:arm64 -f app/Dockerfile_arm64 .
	docker build -t ${DOCKER_USER}/funds-api -f Dockerfile .
	docker build -t ${DOCKER_USER}/funds-api:arm -f Dockerfile_arm .
	docker build -t ${DOCKER_USER}/funds-api:arm64 -f Dockerfile_arm64 .

docker-push:
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASS}
	docker push ${DOCKER_USER}/funds-notifier
	docker push ${DOCKER_USER}/funds-notifier:arm
	docker push ${DOCKER_USER}/funds-notifier:arm64
	docker push ${DOCKER_USER}/funds-front
	docker push ${DOCKER_USER}/funds-front:arm
	docker push ${DOCKER_USER}/funds-front:arm64
	docker push ${DOCKER_USER}/funds-api
	docker push ${DOCKER_USER}/funds-api:arm
	docker push ${DOCKER_USER}/funds-api:arm64
