package fetch

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetBody(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == `/error` {
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer testServer.Close()

	var tests = []struct {
		url     string
		want    []byte
		wantErr error
	}{
		{
			``,
			nil,
			fmt.Errorf(`Error while getting data: Get : unsupported protocol scheme ""`),
		},
		{
			`http://localhost/`,
			nil,
			fmt.Errorf(`Error while getting data: Get http://localhost/: dial tcp [::1]:80: getsockopt: connection refused`),
		},
		{
			testServer.URL + `/error`,
			nil,
			fmt.Errorf(`Error status 400`),
		},
	}

	var failed bool

	for _, test := range tests {
		result, err := GetBody(test.url)

		failed = false

		if err == nil && test.wantErr != nil {
			failed = true
		} else if err != nil && test.wantErr == nil {
			failed = true
		} else if err != nil && err.Error() != test.wantErr.Error() {
			failed = true
		} else if string(result) != string(test.want) {
			failed = true
		}

		if failed {
			t.Errorf(`GetBody(%v) = (%s, %v), want (%s, %v)`, test.url, result, err, test.want, test.wantErr)
		}
	}
}
