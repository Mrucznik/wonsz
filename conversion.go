package wonsz

import "unicode"

// TODO: tests

// ConvertFromCamelCase convert string from camel case format to 2 other formats: dash separated & underscore separated format
func ConvertFromCamelCase(text string) (string, string) {
	separators := map[int]struct{}{}
	longUpper := false
	for i, j := 0, 1; j < len(text); i, j = i+1, j+1 {
		letter := rune(text[i])
		nextLetter := rune(text[j])

		if unicode.IsUpper(letter) && unicode.IsUpper(nextLetter) {
			longUpper = true
			continue
		}

		if unicode.IsUpper(nextLetter) {
			separators[j] = struct{}{}
		} else if longUpper {
			separators[i] = struct{}{}
			longUpper = false
		}
	}

	newLen := len(text) + len(separators)
	dashed := make([]rune, newLen)
	underscored := make([]rune, newLen)
	insrtSeps := 0
	for i, letter := range text {
		idx := i + insrtSeps
		if _, ok := separators[i]; ok {
			dashed[idx] = '-'
			underscored[idx] = '_'
			insrtSeps++
			idx++
		}
		lowered := unicode.ToLower(letter)
		dashed[idx] = lowered
		underscored[idx] = lowered
	}
	return string(dashed), string(underscored)
}
