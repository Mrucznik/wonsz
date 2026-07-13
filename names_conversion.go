package wonsz

import "unicode"

// camelCaseToDashedLowered converts text from
// camelCase naming convention (begin new words with a capital letter except the first word)
// to kebab-case naming convention (separate words with dashes).
func camelCaseToDashedLowered(text string) string {
	return camelCaseToSeparatorsLowered(text, '-')
}

// camelCaseToUnderscoredLowered converts text from
// camelCase naming convention (begin new words with a capital letter except the first word)
// to snake_case naming convention (separate words with underscores).
func camelCaseToUnderscoredLowered(text string) string {
	return camelCaseToSeparatorsLowered(text, '_')
}

// camelCaseToSeparatorsLowered converts text from
// camelCase naming convention (begin new words with a capital letter except the first word)
// to lowercase text with words separated by specified separator.
func camelCaseToSeparatorsLowered(text string, separator rune) string {
	runes := []rune(text)
	separators := getDesiredSeparatorPositions(runes)

	converted := make([]rune, 0, len(runes)+len(separators))
	for i, letter := range runes {
		if _, ok := separators[i]; ok {
			converted = append(converted, separator)
		}
		converted = append(converted, unicode.ToLower(letter))
	}
	return string(converted)
}

// getDesiredSeparatorPositions takes as a parameter text in
// camelCase naming convention (begin new words with a capital letter except the first word)
// and returns a place where two words should be separated (as map keys).
// If a word contains several capital letters in a row, separation will occur before the last capital letter.
// Numbers will also be separated. It skips runes other than letters and numbers.
func getDesiredSeparatorPositions(text []rune) map[int]struct{} {
	separators := map[int]struct{}{}
	upperSequence := false
	numberSequence := false
	for i, j := 0, 1; j < len(text); i, j = i+1, j+1 {
		letter := text[i]
		nextLetter := text[j]

		if !unicode.IsLetter(letter) && !unicode.IsNumber(letter) {
			continue
		}

		if unicode.IsUpper(letter) && unicode.IsUpper(nextLetter) {
			upperSequence = true
			continue
		}

		if unicode.IsNumber(letter) && unicode.IsNumber(nextLetter) {
			numberSequence = true
			continue
		}

		if unicode.IsUpper(nextLetter) || unicode.IsNumber(nextLetter) {
			separators[j] = struct{}{}
		} else if upperSequence {
			separators[i] = struct{}{}
			upperSequence = false
		} else if numberSequence {
			separators[j] = struct{}{}
			numberSequence = false
		}
	}
	return separators
}
