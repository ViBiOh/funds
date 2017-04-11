package morningStar

import (
	"fmt"
	"net/http"
	"testing"
)

func TestGetBody(t *testing.T) {
	var tests = []struct {
		url     string
		httpGet func(string) (*http.Response, error)
		want    []byte
		err     string
	}{
		{
			`test`,
			func(url string) (*http.Response, error) {
				return nil, fmt.Errorf(`Error from test`)
			},
			make([]byte, 0),
			`Error while retrieving data from test: Error from test`,
		},
	}

	for _, test := range tests {
		httpGet = test.httpGet

		body, err := getBody(test.url)
		if err.Error() != test.err || string(body) != string(test.want) {
			t.Errorf("getBody(%v) = (%v, %v), want (%v, %v)", test.url, body, err, test.want, test.err)
		}
	}
}
