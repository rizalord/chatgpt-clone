package helper

import (
	"regexp"
	"strings"
)

func ConvertToSpaced(input string) string {
	// Regular expression to match uppercase letters (except the first one)
	re := regexp.MustCompile(`([a-z])([A-Z])`)
	// Replace matches with a space between the letters
	spaced := re.ReplaceAllString(input, `${1} ${2}`)
	// Convert the first letter to uppercase and the rest to lowercase
	return strings.ToUpper(spaced[:1]) + strings.ToLower(spaced[1:])
}