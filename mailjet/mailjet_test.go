package mailjet

import (
	"testing"
)

func TestPing(t *testing.T) {
	var cases = []struct {
		apiPublicKey string
		want         bool
	}{
		{
			``,
			false,
		},
		{
			`test`,
			true,
		},
	}

	for _, testCase := range cases {
		apiPublicKey = &testCase.apiPublicKey

		if result := Ping(); result != testCase.want {
			t.Errorf(`Ping() = %v, want %v, with apiPublicKey=%v`, result, testCase.want, testCase.apiPublicKey)
		}
	}
}
