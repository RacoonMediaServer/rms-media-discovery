package heuristic

import "strings"

// Normalize returns clear string without non-meanings symbols
func Normalize(text string) string {
	tokens := parse(text)
	result := tokens.String()
	return strings.ToLower(result)
}

func NormalizeWithoutBraces(text string) string {
	tokens := parse(text)
	return tokens.StringWithoutBraces()
}
