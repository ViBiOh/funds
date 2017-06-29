package main

import (
	"github.com/ViBiOh/funds/morningStar"
	"log"
	"net/http"
	"runtime"
)

const port = `1080`

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	http.Handle(`/`, morningStar.Handler{})
	log.Print(`Starting server on port ` + port)
	log.Fatal(http.ListenAndServe(`:`+port, nil))
}
