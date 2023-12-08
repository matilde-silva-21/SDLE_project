package stringStandardizer

import (
	"strings"
	"unicode"
)

func StandardizeString(input string) string {
	// Convert to lowercase
	lowercaseString := strings.ToLower(input)

	// Remove leading and trailing whitespaces
	trimmedString := strings.TrimSpace(lowercaseString)

	// Replace consecutive whitespaces with a single underscore
	replacedWhitespace := ReplaceConsecutiveWhitespace(trimmedString, '_')

	return replacedWhitespace
}

func ReplaceConsecutiveWhitespace(input string, replacement rune) string {
	var result strings.Builder
	prevWhitespace := false

	for _, char := range input {
		if unicode.IsSpace(char) {
			// If it's the first whitespace or consecutive whitespaces, append the replacement
			if !prevWhitespace {
				result.WriteRune(replacement)
			}
			prevWhitespace = true
		} else {
			// If it's not a whitespace, append the character
			result.WriteRune(char)
			prevWhitespace = false
		}
	}

	return result.String()
}
