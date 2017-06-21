default: deps lint coverage build

deps:
	go get -u github.com/golang/lint/golint

lint:
	golint ./...
	go vet ./...

coverage:
	./test/coverage

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo server.go
