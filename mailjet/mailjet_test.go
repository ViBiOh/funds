package mailjet

import (
	"testing"
)

func TestPing(t *testing.T) {
	var tests = []struct {
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

	for _, test := range tests {
		apiPublicKey = &test.apiPublicKey

		if result := Ping(); result != test.want {
			t.Errorf(`Ping() = %v, want %v, with apiPublicKey=%v`, result, test.want, test.apiPublicKey)
		}
	}
}
