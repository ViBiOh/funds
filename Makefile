default: lint vet coverage build

lint:
	go get -u github.com/golang/lint/golint
	golint ./...

vet:
	go vet ./...

coverage:
	./test/coverage

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo server.go
