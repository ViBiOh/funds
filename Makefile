default: deps dev

dev: format lint tst bench build

deps:
	go get -u github.com/golang/lint/golint
	go get -u github.com/lib/pq
	go get -u github.com/NYTimes/gziphandler
	go get -u github.com/tdewolff/minify
	go get -u github.com/ViBiOh/alcotest/alcotest
	go get -u github.com/ViBiOh/httputils
	go get -u github.com/ViBiOh/httputils/cors
	go get -u github.com/ViBiOh/httputils/db
	go get -u github.com/ViBiOh/httputils/owasp
	go get -u github.com/ViBiOh/httputils/prometheus
	go get -u github.com/ViBiOh/httputils/rate
	go get -u github.com/ViBiOh/httputils/tools
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
