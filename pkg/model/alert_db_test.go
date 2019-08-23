package model

import (
	"testing"

	"github.com/ViBiOh/httputils/pkg/errors"
)

func TestSaveAlert(t *testing.T) {
	var cases = []struct {
		intention string
		input     *Alert
		wantErr   error
	}{
		{
			"simple",
			nil,
			errors.New("cannot save nil"),
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.intention, func(t *testing.T) {
			app := app{}

			err := app.SaveAlert(testCase.input, nil)

			failed := false

			if err == nil && testCase.wantErr != nil {
				failed = true
			} else if err != nil && testCase.wantErr == nil {
				failed = true
			} else if err != nil && err.Error() != testCase.wantErr.Error() {
				failed = true
			}

			if failed {
				t.Errorf("SaveAlert(%#v) = (%#v), want (%#v)", testCase.input, err, testCase.wantErr)
			}
		})
	}
}
