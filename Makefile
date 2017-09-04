default: deps format lint tst build

deps:
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/golang/lint/golint
	go get -u github.com/ViBiOh/httputils
	go get -u github.com/ViBiOh/httputils/cors
	go get -u github.com/ViBiOh/httputils/owasp
	go get -u github.com/ViBiOh/httputils/prometheus
	go get -u github.com/ViBiOh/httputils/tools
	go get -u github.com/ViBiOh/alcotest/alcotest
	go get -u github.com/lib/pq
	go get -u github.com/tdewolff/minify

format:
	goimports -w **/*.go
	gofmt -s -w **/*.go

lint:
	golint ./...
	go vet ./...

tst:
	script/coverage

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo -o bin/funds api/api.go
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo -o bin/notifier alert/alert.go
