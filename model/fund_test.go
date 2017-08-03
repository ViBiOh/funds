package model

import (
	"testing"
)

func TestGetID(t *testing.T) {
	var tests = []struct {
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

	for _, test := range tests {
		result := test.instance.GetID()
		if result != test.want {
			t.Errorf(`GetID() of %v = %v, want %v`, test.instance, result, test.want)
		}
	}
}

func TestComputeScore(t *testing.T) {
	var tests = []struct {
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

	for _, test := range tests {
		test.instance.ComputeScore()
		if test.instance.Score != test.want {
			t.Errorf(`ComputeScore() of %v = %v, want %v`, test.instance, test.instance.Score, test.want)
		}
	}
}
