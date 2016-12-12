default: back

back:
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo go/server.go
	go test -run ./go/
