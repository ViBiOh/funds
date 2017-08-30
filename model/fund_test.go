package model

import (
	"testing"
)

func TestGetID(t *testing.T) {
	var cases = []struct {
		instance Fund
		want     string
	}{
		{
			Fund{},
			``,
		},
		{
			Fund{ID: `test`},
			`test`,
		},
	}

	for _, testCase := range cases {
		result := testCase.instance.GetID()
		if result != testCase.want {
			t.Errorf(`GetID() of %v = %v, want %v`, testCase.instance, result, testCase.want)
		}
	}
}

func TestComputeScore(t *testing.T) {
	var cases = []struct {
		instance *Fund
		want     float64
	}{
		{
			&Fund{},
			0.0,
		},
		{
			&Fund{OneMonth: 1 / 0.25, ThreeMonths: 1 / 0.3, SixMonths: 1 / 0.25, OneYear: 1 / 0.2, VolThreeYears: 1 / 0.1},
			3.0,
		},
	}

	for _, testCase := range cases {
		testCase.instance.ComputeScore()
		if testCase.instance.Score != testCase.want {
			t.Errorf(`ComputeScore() of %v = %v, want %v`, testCase.instance, testCase.instance.Score, testCase.want)
		}
	}
}
