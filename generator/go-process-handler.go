package generator

import (
	"fmt"
	"regexp"
)

func processQuery(query string) (string, map[string]string) {
	// Regex patterns
	extrapolateRegex := regexp.MustCompile(`\{([a-zA-Z0-9_]+)\}\$`) // Matches `{XXX}$`
	insertRegex := regexp.MustCompile(`\$\{([a-zA-Z0-9_]+)\}`)      // Matches `${XXX}`

	// Step 1: Replace {XXX}$ → XXX
	query = extrapolateRegex.ReplaceAllString(query, "$1")

	// Step 2: Replace ${XXX} with $1, $2, etc., and track numbers
	placeholderMap := make(map[string]string)
	counter := 1

	query = insertRegex.ReplaceAllStringFunc(query, func(match string) string {
		// Extract key inside ${XXX}
		key := match[2 : len(match)-1] // Removes ${ and }
		if _, exists := placeholderMap[key]; !exists {
			placeholderMap[key] = fmt.Sprintf("$%d", counter)
			counter++
		}
		return placeholderMap[key]
	})

	return query, placeholderMap
}

func processHandler(handler string, queryParams map[string]string) string {
	// Replace `${}` with `w, r`
	handler = regexp.MustCompile(`\$\{\}`).ReplaceAllString(handler, "w, r")

	// Replace `${example}` (or any query param key) with its value
	paramRegex := regexp.MustCompile(`\$\{([a-zA-Z0-9_]+)\}`)
	handler = paramRegex.ReplaceAllStringFunc(handler, func(match string) string {
		key := match[2 : len(match)-1] // Extracts key inside `${example}`
		if val, exists := queryParams[key]; exists {
			return val
		}
		return match // Keep unchanged if no matching param found
	})

	return handler
}
