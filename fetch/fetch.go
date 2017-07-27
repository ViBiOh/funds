package fetch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

var httpClient = http.Client{Timeout: 30 * time.Second}

func readBody(body io.ReadCloser) ([]byte, error) {
	defer body.Close()
	return ioutil.ReadAll(body)
}

// GetBody return body of given URL or error if something goes wrong
func GetBody(url string) ([]byte, error) {
	request, err := http.NewRequest(`GET`, url, nil)
	if err != nil {
		return nil, fmt.Errorf(`Unable to prepare request for url %s : %v`, url, err)
	}

	response, err := httpClient.Do(request)
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

// PostJSONBody post given interface to URL with optional credential supplied
func PostJSONBody(url string, body interface{}, user string, pass string) ([]byte, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(`POST`, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	request.Header.Add(`Content-Type`, `application/json`)
	if user != `` {
		request.SetBasicAuth(user, pass)
	}

	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	responseContent, err := readBody(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf(`Got status %d while sending mail %s`, response.StatusCode, string(responseContent))
	}

	return responseContent, nil
}
