default: deps lint tst build

deps:
	go get -u github.com/golang/lint/golint

lint:
	golint ./...
	go vet ./...

tst:
	./test/coverage

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo funds.go
