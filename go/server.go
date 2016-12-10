package main

import (
	"github.com/ViBiOh/funds-ob/go/morningStar"
	"log"
	"net/http"
	"runtime"
	"strconv"
)

const port = `1080`

func main() {
	numCpu := runtime.NumCPU()
	runtime.GOMAXPROCS(numCpu)
	log.Print(`MaxProc setted to ` + strconv.Itoa(numCpu))

	http.Handle(`/`, morningStar.Handler{})
	log.Print(`Starting server on port ` + port)
	log.Fatal(http.ListenAndServe(`:`+port, nil))
}
