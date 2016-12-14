default: test vet build

test:
	go test ./...

test:
	go vet ./...

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo server.go
