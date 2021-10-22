package main

import "unicode"

func CamelCaseToDashed(text string) string {
	return CamelCaseToSeparators(text, '-')
}

func CamelCaseToUnderscored(text string) string {
	return CamelCaseToSeparators(text, '_')
}

func CamelCaseToSeparators(text string, separator rune) string {
	separators := getDesiredSeparatorPositions(text)

	newLen := len(text) + len(separators)
	converted := make([]rune, newLen)
	insertedSeparators := 0
	for i, letter := range text {
		idx := i + insertedSeparators
		if _, ok := separators[i]; ok {
			converted[idx] = separator
			insertedSeparators++
			idx++
		}
		converted[idx] = letter
	}
	return string(converted)
}

func getDesiredSeparatorPositions(text string) map[int]struct{} {
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
	return separators
}
