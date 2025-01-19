package generator

import (
	"bufio"
	"fmt"
	"master-gen/internal/parser"
	"os"
	"regexp"
	"strings"
)

func generateHandler(op parser.Operation, structParams map[string]map[string]string, singleParams map[string]string) string {
	var goCode string

	goCode = genQuery(op, structParams, singleParams)

	goCode += `
	h.JSON.Success(w, res)
}	
`
	return goCode
}

func genQuery(op parser.Operation, structParams map[string]map[string]string, singleParams map[string]string) string {
	goCode := ""
	goCode += `
func (h *Handler) ` + op.Query + `(w http.ResponseWriter, r *http.Request) {`

	goCode += generateBody(op.Query, structParams, singleParams)

	paramStr := genParamStr(op.Query, structParams, singleParams)

	goCode += generateQueryBinding(op.Query, paramStr)

	return goCode
}

func genParamStr(query string, structParams map[string]map[string]string, singleParams map[string]string) string {
	// Check for struct-based parameters
	if _, structExists := structParams[query+"Params"]; structExists {
		return ", body"
	}
	// Check for single parameter
	if singleParam, singleExists := singleParams[query]; singleExists {
		parts := strings.Split(singleParam, " ")
		paramName := parts[0]
		return ", " + paramName
	}

	// No parameters
	return ""
}

func isExecQuery(queryName string) bool {
	file, err := os.Open("./queries.sql")
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "-- name:") {
			parts := strings.Fields(line)
			if len(parts) >= 3 && parts[2] == queryName {
				return parts[3] == ":exec"
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return false
	}

	return false
}

func generateQueryBinding(query, paramStr string) string {
	insert := "res, err"
	if isExecQuery(query) {
		insert = `res := make(map[string]string)
	err`
	}

	return fmt.Sprintf(`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	%s := h.Queries.%s(ctx%s)
	if err != nil {
		h.JSON.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	`, insert, query, paramStr)
}

func generateQueryParams(query string, structParams map[string]map[string]string, singleParams map[string]string) string {
	paramStruct := query + "Params"
	fields, structExists := structParams[paramStruct]
	singleParam, singleExists := singleParams[query]

	if structExists {
		// Handle multi-field query params
		goCode := `
	queryParams := queries.` + paramStruct + `{`
		for field, fieldType := range fields {
			goCode += `
		` + field + `: r.URL.Query().Get("` + strings.ToLower(field) + `"), // Type: ` + fieldType
		}
		goCode += `
	}
	`
		return goCode
	} else if singleExists {
		// Handle single-parameter case
		parts := strings.Split(singleParam, " ")
		paramName := parts[0]
		paramType := parts[1]

		goCode := `
	var ` + paramName + ` ` + paramType + `
	` + paramName + ` = r.URL.Query().Get("` + strings.ToLower(paramName) + `") // Type: ` + paramType + `
	`
		return goCode
	}

	return "// No Params Found"
}

func generateBody(query string, bodies map[string]map[string]string, singleParams map[string]string) string {
	bodyStruct := query + "Params"
	_, structExists := bodies[bodyStruct]
	singleParam, singleExists := singleParams[query]

	if structExists {
		// Handle multi-field body structs
		goCode := `
		body := queries.` + bodyStruct + `{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			h.JSON.Error(w, http.StatusBadRequest, "Failed to decode body")
			return
		}
		`
		return goCode
	} else if singleExists {
		// Handle single parameter case
		parts := strings.Split(singleParam, " ")
		paramName := parts[0]
		paramType := parts[1]
		goCode := `
		var ` + paramName + ` ` + paramType + `
		if err := json.NewDecoder(r.Body).Decode(&` + paramName + `); err != nil {
			h.JSON.Error(w, http.StatusBadRequest, "Failed to decode body")
			return
		}
		`
		return goCode
	}

	return "// No Body Found"
}

func parseSqlcFile(filePath string) (map[string]map[string]string, map[string]string, error) {
	structParams := make(map[string]map[string]string)
	singleParams := make(map[string]string)

	var currentStruct string
	var insideStruct bool

	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Detect struct definition
		if strings.HasPrefix(line, "type") && strings.Contains(line, "struct {") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				currentStruct = parts[1]
				structParams[currentStruct] = make(map[string]string)
				insideStruct = true
			}
			continue
		}

		// Parse struct fields
		if insideStruct {
			if line == "}" {
				insideStruct = false
				currentStruct = ""
				continue
			}
			parts := strings.Fields(line)
			if len(parts) == 2 && currentStruct != "" {
				fieldName := parts[0]
				fieldType := parts[1]
				structParams[currentStruct][fieldName] = fieldType
			}
			continue
		}

		// Detect function signature
		if strings.HasPrefix(line, "func (q *Queries)") {
			start := strings.Index(line, "(")
			end := strings.LastIndex(line, ")")
			if start != -1 && end != -1 {
				funcSig := line[start+1 : end]
				params := strings.Split(funcSig, ",")

				// Skip ctx parameter
				if len(params) == 2 {
					paramParts := strings.Fields(params[1])
					if len(paramParts) == 2 {
						paramName := paramParts[0]
						paramType := paramParts[1]
						methodName := extractMethodName(line)

						if paramName == "arg" {
							if methodName != "" {
								structParams[methodName+"Params"] = make(map[string]string)
							}
						} else {
							if methodName != "" {
								singleParams[methodName] = fmt.Sprintf("%s %s", paramName, paramType)
							}
						}
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return structParams, singleParams, nil
}

func extractMethodName(line string) string {
	// Regex to match and extract the method name in the function signature
	methodRegex := regexp.MustCompile(`func\s+\(q\s+\*Queries\)\s+(\w+)\s*\(ctx\s+context\.Context,\s*(\w+)\s+([\w\.\[\]]+)\)`)
	matches := methodRegex.FindStringSubmatch(line)

	if len(matches) > 1 {
		return matches[1] // The method name is the first capture group
	}
	return ""
}
