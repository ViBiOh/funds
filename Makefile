default: back

back:
	go test ./go/*/
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo go/server.go
