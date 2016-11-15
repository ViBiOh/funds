package morningStarWs

import (
	"../morningStar"
	"encoding/json"
	"golang.org/x/net/websocket"
)

func handleError(ws *websocket.Conn, err error) bool {
	if err != nil {
		ws.Write([]byte(err.Error()))
		return true
	}
	return false
}

func Handler(ws *websocket.Conn) {
	msg := make([]byte, 512)

	n, err := ws.Read(msg)
	if handleError(ws, err) {
		return
	}

	performance, err := morningStar.SinglePerformance(string(msg[:n]))
	if handleError(ws, err) {
		return
	}

	objJson, err := json.Marshal(performance)
	if handleError(ws, err) {
		return
	}

	_, err = ws.Write(objJson)
	if handleError(ws, err) {
		return
	}
}
