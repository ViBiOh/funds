package model

import (
	"regexp"
	"testing"
)

func TestCleanID(t *testing.T) {
	var tests = []struct {
		fundID []byte
		want   string
	}{
		{},
		{
			[]byte(`aZeRtY`),
			`azerty`,
		},
	}

	for _, test := range tests {
		if got := cleanID(test.fundID); string(got) != test.want {
			t.Errorf(`cleanID(%v) = %v, want %v`, test.fundID, got, test.want)
		}
	}
}

func TestExtractLabel(t *testing.T) {
	var tests = []struct {
		extract      *regexp.Regexp
		body         []byte
		defaultValue []byte
		want         string
	}{
		{},
		{
			regexp.MustCompile(`id:(\S+)`),
			[]byte(`id:12345`),
			[]byte(``),
			`12345`,
		},
		{
			regexp.MustCompile(`id:\S+`),
			[]byte(`id:12345`),
			[]byte(``),
			``,
		},
		{
			regexp.MustCompile(`label:(\S+)`),
			[]byte(`I'm looking to extract an id:12345 in this body`),
			[]byte(`notFound`),
			`notFound`,
		},
		{
			regexp.MustCompile(`label:(\S+)`),
			[]byte(`label:Alice&amp;Bob`),
			[]byte(`notFound`),
			`Alice&Bob`,
		},
	}

	for _, test := range tests {
		if got := extractLabel(test.extract, test.body, test.defaultValue); string(got) != test.want {
			t.Errorf(`extractLabel(%v, %v, %v) = %v, want %v`, test.extract, test.body, test.defaultValue, got, test.want)
		}
	}
}

func TestExtractPerformance(t *testing.T) {
	var tests = []struct {
		extract *regexp.Regexp
		body    []byte
		want    float64
	}{
		{
			regexp.MustCompile(`ISIN.:(\S+)`),
			[]byte(`ISIN :3.14%`),
			3.14,
		},
		{
			regexp.MustCompile(`ISIN.:(\S+)`),
			[]byte(`ISIN :-.07%`),
			-0.07,
		},
		{
			regexp.MustCompile(`ISIN.:(\S+)`),
			[]byte(`ISIN :notValid`),
			0.0,
		},
	}

	for _, test := range tests {
		if got := extractPerformance(test.extract, test.body); got != test.want {
			t.Errorf(`extractPerformance(%v, %v) = %v, want %v`, test.extract, test.body, got, test.want)
		}
	}
}
