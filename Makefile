default: back

back:
	go get golang.org/x/net/websocket
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo src/server.go
