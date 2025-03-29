package utils

import (
	"strings"
	"unicode"
)

// VariableName
func ToPascalCase(input string) string {
	splittedWords := splitWords(input)

	for i := range splittedWords {
		splittedWords[i] = strings.Title(splittedWords[i])
	}

	return strings.Join(splittedWords, "")
}

// variableName
func ToCamelCase(input string) string {
	// Change the first word into lowerCase
	splittedWords := splitWords(input)

	for i := range splittedWords {
		if i == 0 {
			splittedWords[i] = strings.ToLower(splittedWords[i])
		} else {
			splittedWords[i] = strings.Title(splittedWords[i])
		}
	}

	return strings.Join(splittedWords, "")
}

// variable_name
func ToSnakeCase(input string) string {
	return strings.ToLower(strings.Join(splitWords(input), "_"))
}

// variable-name
func ToKebabCase(input string) string {
	return strings.ToLower(strings.Join(splitWords(input), "-"))
}

func splitWords(value string) []string {
	// Store all words within value
	var words []string
	// Store all characters of single word
	var currWord []rune

	// Loop the `value`
	for _, ch := range value {
		// Check if `ch` is Letter or digit
		// Then add `ch` to the `currWord`
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) {
			currWord = append(currWord, ch)
		} else if len(currWord) > 0 {
			// else if `ch` is not Letter or digit, then it entered new word
			// hence appending `word` array with `currWord`
			words = append(words, string(currWord))
			// Then reset the `currWord` to accept new word
			currWord = nil
		}
	}

	// Handle if `currWord` doesn't get reset inside the loop
	if len(currWord) > 0 {
		words = append(words, string(currWord))
	}

	return words
}
