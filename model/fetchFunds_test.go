package model

import (
	"regexp"
	"testing"
)

func TestCleanID(t *testing.T) {
	var cases = []struct {
		fundID []byte
		want   string
	}{
		{},
		{
			[]byte(`aZeRtY`),
			`azerty`,
		},
	}

	for _, testCase := range cases {
		if got := cleanID(testCase.fundID); string(got) != testCase.want {
			t.Errorf(`cleanID(%v) = %v, want %v`, testCase.fundID, got, testCase.want)
		}
	}
}

func TestExtractLabel(t *testing.T) {
	var cases = []struct {
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

	for _, testCase := range cases {
		if got := extractLabel(testCase.extract, testCase.body, testCase.defaultValue); string(got) != testCase.want {
			t.Errorf(`extractLabel(%v, %v, %v) = %v, want %v`, testCase.extract, testCase.body, testCase.defaultValue, got, testCase.want)
		}
	}
}

func TestExtractPerformance(t *testing.T) {
	var cases = []struct {
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

	for _, testCase := range cases {
		if got := extractPerformance(testCase.extract, testCase.body); got != testCase.want {
			t.Errorf(`extractPerformance(%v, %v) = %v, want %v`, testCase.extract, testCase.body, got, testCase.want)
		}
	}
}
