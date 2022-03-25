package model

import (
	"regexp"
	"testing"
)

func TestCleanID(t *testing.T) {
	cases := map[string]struct {
		fundID []byte
		want   string
	}{
		"simple": {
			fundID: []byte("aZeRtY"),
			want:   "azerty",
		},
	}

	for intention, tc := range cases {
		t.Run(intention, func(t *testing.T) {
			if got := cleanID(tc.fundID); got != tc.want {
				t.Errorf("cleanID(%v) = %v, want %v", tc.fundID, got, tc.want)
			}
		})
	}
}

func TestExtractLabel(t *testing.T) {
	cases := map[string]struct {
		extract      *regexp.Regexp
		body         []byte
		defaultValue []byte
		want         string
	}{
		"simple": {
			regexp.MustCompile(`id:(\S+)`),
			[]byte("id:12345"),
			[]byte(""),
			"12345",
		},
		"no match": {
			regexp.MustCompile(`id:\S+`),
			[]byte("id:12345"),
			[]byte(""),
			"",
		},
		"not found": {
			regexp.MustCompile(`label:(\S+)`),
			[]byte("I'm looking to extract an id:12345 in this body"),
			[]byte("notFound"),
			"notFound",
		},
		"found": {
			regexp.MustCompile(`label:(\S+)`),
			[]byte("label:Alice&amp;Bob"),
			[]byte("notFound"),
			"Alice&Bob",
		},
	}

	for intention, tc := range cases {
		t.Run(intention, func(t *testing.T) {
			if got := extractLabel(tc.extract, tc.body, tc.defaultValue); string(got) != tc.want {
				t.Errorf("extractLabel(%v, %v, %v) = %v, want %v", tc.extract, tc.body, tc.defaultValue, got, tc.want)
			}
		})
	}
}

func TestExtractPerformance(t *testing.T) {
	cases := map[string]struct {
		extract *regexp.Regexp
		body    []byte
		want    float64
	}{
		"positive": {
			regexp.MustCompile(`ISIN.:(\S+)`),
			[]byte("ISIN :3.14%"),
			3.14,
		},
		"negative": {
			regexp.MustCompile(`ISIN.:(\S+)`),
			[]byte("ISIN :-.07%"),
			-0.07,
		},
		"invalid": {
			regexp.MustCompile(`ISIN.:(\S+)`),
			[]byte("ISIN :notValid"),
			0.0,
		},
	}

	for intention, tc := range cases {
		t.Run(intention, func(t *testing.T) {
			if got := extractPerformance(tc.extract, tc.body); got != tc.want {
				t.Errorf("extractPerformance(%v, %v) = %v, want %v", tc.extract, tc.body, got, tc.want)
			}
		})
	}
}
