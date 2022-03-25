package model

import (
	"context"
	"errors"
	"testing"
)

func TestSaveAlert(t *testing.T) {
	cases := map[string]struct {
		input   *Alert
		wantErr error
	}{
		"simple": {
			nil,
			errors.New("cannot save nil"),
		},
	}

	for intention, tc := range cases {
		t.Run(intention, func(t *testing.T) {
			app := App{}

			err := app.SaveAlert(context.Background(), tc.input)

			failed := false

			if err == nil && tc.wantErr != nil {
				failed = true
			} else if err != nil && tc.wantErr == nil {
				failed = true
			} else if err != nil && err.Error() != tc.wantErr.Error() {
				failed = true
			}

			if failed {
				t.Errorf("SaveAlert() = `%s`, want `%s`", err, tc.wantErr)
			}
		})
	}
}
