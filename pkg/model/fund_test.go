package model

import (
	"testing"
)

func TestGetID(t *testing.T) {
	cases := map[string]struct {
		instance Fund
		want     string
	}{
		"empty": {
			Fund{},
			"",
		},
		"value": {
			Fund{ID: "test"},
			"test",
		},
	}

	for intention, tc := range cases {
		t.Run(intention, func(t *testing.T) {
			result := tc.instance.GetID()
			if result != tc.want {
				t.Errorf("GetID() of %v = %v, want %v", tc.instance, result, tc.want)
			}
		})
	}
}

func TestComputeScore(t *testing.T) {
	cases := map[string]struct {
		instance *Fund
		want     float64
	}{
		"empty": {
			&Fund{},
			0.0,
		},
		"value": {
			&Fund{OneMonth: 1 / 0.25, ThreeMonths: 1 / 0.3, SixMonths: 1 / 0.25, OneYear: 1 / 0.2, VolThreeYears: 1 / 0.1},
			3.0,
		},
	}

	for intention, tc := range cases {
		t.Run(intention, func(t *testing.T) {
			tc.instance.ComputeScore()
			if tc.instance.Score != tc.want {
				t.Errorf("ComputeScore() of %v = %v, want %v", tc.instance, tc.instance.Score, tc.want)
			}
		})
	}
}
