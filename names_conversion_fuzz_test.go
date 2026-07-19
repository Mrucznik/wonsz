package wonsz

import (
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

func FuzzCamelCaseToUnderscoredLowered(f *testing.F) {
	f.Add("simpleTestingName")
	f.Add("ŁadnyWąż")
	f.Add("userWith5BANSButNoMoreTHAN100Ok1200ok")
	f.Add("Non-ConventionalStringMyID")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		got := camelCaseToUnderscoredLowered(input)

		// The pre-rune-conversion implementation introduced NUL runes for
		// multibyte input; the function must never add NULs on its own.
		if !strings.ContainsRune(input, 0) && strings.ContainsRune(got, 0) {
			t.Errorf("output contains introduced NUL rune: %q -> %q", input, got)
		}
		for _, r := range got {
			// Some uppercase runes (e.g. U+03D4) have no lowercase mapping;
			// the invariant is that nothing lowercasable stays uppercase.
			if unicode.IsUpper(r) && unicode.ToLower(r) != r {
				t.Errorf("output contains lowercasable uppercase rune %q: %q -> %q", r, input, got)
				break
			}
		}
		if utf8.ValidString(input) && !utf8.ValidString(got) {
			t.Errorf("valid UTF-8 input produced invalid UTF-8 output: %q -> %q", input, got)
		}
	})
}
