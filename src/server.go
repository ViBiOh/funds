package main

import (
	"./morningStar"
	"./morningStarWs"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

const port = `1080`

func main() {
	http.Handle(`/ws/`, websocket.Handler(morningStarWs.Handler))
	http.HandleFunc(`/`, morningStar.Handler)

	log.Print(`Starting server on port ` + port)
	log.Fatal(http.ListenAndServe(`:`+port, nil))
}
