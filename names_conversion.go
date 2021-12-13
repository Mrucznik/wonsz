package wonsz

import "unicode"

// TODO: Here some description
func camelCaseToDashedLowered(text string) string {
	return camelCaseToSeparatorsLowered(text, '-')
}

// TODO: Here some description
func camelCaseToUnderscoredLowered(text string) string {
	return camelCaseToSeparatorsLowered(text, '_')
}

// TODO: Here some description
func camelCaseToSeparatorsLowered(text string, separator rune) string {
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
		lowered := unicode.ToLower(letter)
		converted[idx] = lowered
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
