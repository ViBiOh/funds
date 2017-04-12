package morningStar

import (
	"fmt"
	"net/http"
	"testing"
)

type FakeReaderCloser struct {
	data []byte
	err  error
}

func (o FakeReaderCloser) Read([]byte) (int, error) {
	if o.err != nil {
		return 0, o.err
	}
	return 0, nil
}

func (FakeReaderCloser) Close() error {
	return nil
}

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
		{
			`test`,
			func(url string) (*http.Response, error) {
				return &http.Response{StatusCode: 401}, nil
			},
			make([]byte, 0),
			`Got error 401 while getting test`,
		},
		{
			`test`,
			func(url string) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       FakeReaderCloser{nil, fmt.Errorf(`Error from test`)},
				}, nil
			},
			make([]byte, 0),
			`Error while reading body of test: Error from test`,
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
