package model

import (
	"context"
	"errors"
	"testing"
)

func TestSaveAlert(t *testing.T) {
	cases := []struct {
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
			app := App{}

			err := app.SaveAlert(context.Background(), testCase.input)

			failed := false

			if err == nil && testCase.wantErr != nil {
				failed = true
			} else if err != nil && testCase.wantErr == nil {
				failed = true
			} else if err != nil && err.Error() != testCase.wantErr.Error() {
				failed = true
			}

			if failed {
				t.Errorf("SaveAlert() = `%s`, want `%s`", err, testCase.wantErr)
			}
		})
	}
}
