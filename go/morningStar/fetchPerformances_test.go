package morningStar

import (
	"regexp"
	"testing"
)

func TestExtractLabel(t *testing.T) {
	extract := regexp.MustCompile(`id:(\S+)`)
	body := []byte(`I'm looking to extract an id:12345 in this body`)
	defaultValue := []byte(``)

	if got := extractLabel(extract, body, defaultValue); string(got) != `12345` {
		t.Errorf("extractLabel(%q, %q, %q) = %v", extract, body, defaultValue, got)
	}
}
