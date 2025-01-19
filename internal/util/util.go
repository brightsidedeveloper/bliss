package util

import "strings"

func Capitalize(s string) string {
	return strings.ToUpper(string(s[0])) + s[1:]
}

func Lower(s string) string {
	return strings.ToLower(s)
}

func ConvertToCamelCase(input string) string {
	words := strings.Split(input, "_")
	for i := range words {
		words[i] = strings.Title(words[i])
	}
	return strings.Join(words, "")
}
