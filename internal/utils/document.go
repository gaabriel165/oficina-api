package utils

import "regexp"

var nonDigitRegex = regexp.MustCompile(`\D`)

func CleanNumericString(value string) string {
	return nonDigitRegex.ReplaceAllString(value, "")
}
