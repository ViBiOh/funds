package morningStar

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
)

type FakeReaderCloser struct {
	reader io.Reader
	err    error
}

func (o FakeReaderCloser) Read(p []byte) (int, error) {
	if o.err != nil {
		return 0, o.err
	}
	return o.reader.Read(p)
}

func (FakeReaderCloser) Close() error {
	return nil
}

func fakeGet(err error, statusCode int, body []byte, readerErr error) func(string) (*http.Response, error) {
	return func(string) (*http.Response, error) {
		return &http.Response{
			StatusCode: statusCode,
			Body:       FakeReaderCloser{bytes.NewReader(body), readerErr},
		}, err
	}
}

func TestGetBody(t *testing.T) {
	var tests = []struct {
		url        string
		err        error
		statusCode int
		body       []byte
		readerErr  error
		want       []byte
		wantErr    string
	}{
		{
			url:     `test`,
			err:     fmt.Errorf(`Error from get`),
			wantErr: `Error while retrieving data from test: Error from get`,
		},
		{
			url:        `test`,
			statusCode: 401,
			wantErr:    `Got error 401 while getting test`,
		},
		{
			url:        `test`,
			statusCode: 200,
			readerErr:  fmt.Errorf(`Error from reader`),
			wantErr:    `Error while reading body of test: Error from reader`,
		},
		{
			url:        `test`,
			statusCode: 200,
			body:       []byte(`Body retrieved from fetch`),
			want:       []byte(`Body retrieved from fetch`),
		},
	}

	for _, test := range tests {
		httpGet = fakeGet(test.err, test.statusCode, test.body, test.readerErr)

		body, err := getBody(test.url)
		if (err == nil && test.wantErr != ``) || (err != nil && err.Error() != test.wantErr) || string(body) != string(test.want) {
			t.Errorf("getBody(%v) = (%v, %v), want (%v, %v)", test.url, body, err, test.want, test.wantErr)
		}
	}
}
