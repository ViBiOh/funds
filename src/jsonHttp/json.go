package jsonHttp

import (
	"encoding/json"
	"net/http"
)

func ResponseJson(w http.ResponseWriter, obj interface{}) {
	objJson, err := json.Marshal(obj)

	if err == nil {
		w.Header().Set(`Content-Type`, `application/json`)
		w.Header().Set(`Cache-Control`, `no-cache`)
		w.Header().Set(`Access-Control-Allow-Origin`, `*`)
		w.Write(objJson)
	} else {
		http.Error(w, `Error while marshalling JSON response`, 500)
	}
}
