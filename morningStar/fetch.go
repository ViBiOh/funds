package morningStar

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

var httpGet = http.Get

func readBody(body io.ReadCloser) ([]byte, error) {
	defer body.Close()
	return ioutil.ReadAll(body)
}

func getBody(url string) ([]byte, error) {
	response, err := httpGet(url)
	if err != nil {
		return nil, fmt.Errorf(`Error while retrieving data from %s: %v`, url, err)
	}

	if response.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf(`Got error %d while getting %s`, response.StatusCode, url)
	}

	body, err := readBody(response.Body)
	if err != nil {
		return nil, fmt.Errorf(`Error while reading body of %s: %v`, url, err)
	}

	return body, nil
}
