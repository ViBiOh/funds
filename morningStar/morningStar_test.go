package morningStar

import (
	"testing"
)

func TestComputeScore(t *testing.T) {
	var tests = []struct {
		instance *performance
		want     float64
	}{
		{
			&performance{},
			0.0,
		},
		{
			&performance{OneMonth: 1 / 0.25, ThreeMonths: 1 / 0.3, SixMonths: 1 / 0.25, OneYear: 1 / 0.2, VolThreeYears: 1 / 0.1},
			3.0,
		},
	}

	for _, test := range tests {
		test.instance.computeScore()
		if test.instance.Score != test.want {
			t.Errorf("computeScore() with %v = %v, want %v", test.instance, test.instance.Score, test.want)
		}
	}
}
