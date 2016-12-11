package morningStar

import (
	"regexp"
	"testing"
)

func TestExtractLabel(t *testing.T) {
	var tests = []struct {
		extract      *regexp.Regexp
		body         []byte
		defaultValue []byte
		want         string
	}{
		{
			regexp.MustCompile(`id:(\S+)`),
			[]byte(`I'm looking to extract an id:12345 in this body`),
			[]byte(``),
			`12345`,
		},
	}

	for _, test := range tests {
		if got := extractLabel(test.extract, test.body, test.defaultValue); string(got) != test.want {
			t.Errorf("extractLabel(%q, %q, %q) = %v", test.extract, test.body, test.defaultValue, got)
		}
	}

}
