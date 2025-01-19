package generator

import (
	"bufio"
	"fmt"
	"master-gen/internal/parser"
	"os"
	"regexp"
	"strings"
)

func generateHandler(op parser.Operation, params SqlcParams, types SqlcParams) string {
	var goCode string

	goCode = genQuery(op, params, types)

	goCode += `
	h.JSON.Success(w, res)
}	
`
	return goCode
}

func genQuery(op parser.Operation, params SqlcParams, types SqlcParams) string {
	goCode := ""
	name := op.Query
	if name == "" {
		name = op.Handler
	}
	goCode += `
func (h *Handler) ` + name + `(w http.ResponseWriter, r *http.Request) {`

	if op.Query != "" {
		goCode += generateBody(op.Query, params, false)
		paramStr := genParamStr(op.Query, params)
		goCode += generateQueryBinding(op.Query, paramStr)
	} else {
		goCode += generateBody(op.Handler, types, true)
		paramStr := genParamStr(op.Handler, types)
		goCode += generateHandlerBinding(op.Handler, paramStr)
	}

	return goCode
}

func genParamStr(query string, params SqlcParams) string {
	if _, structExists := params.structured[query+"Params"]; structExists {
		return ", body"
	}

	if singleParam, singleExists := params.single[query]; singleExists {
		parts := strings.Split(singleParam, " ")
		paramName := parts[0]
		return ", " + paramName
	}

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
func generateHandlerBinding(handler, paramStr string) string {
	insert := "res, err"

	return fmt.Sprintf(`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	%s := h.Injections.%s(ctx%s)
	if err != nil {
		h.JSON.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	`, insert, handler, paramStr)
}

func generateBody(query string, params SqlcParams, isHandler bool) string {
	bodyStruct := query + "Params"
	_, structExists := params.structured[bodyStruct]
	singleParam, singleExists := params.single[query]

	from := "queries"
	if isHandler {
		from = "injections"
	}

	if structExists {
		goCode := fmt.Sprintf(`
	body := %s.%s{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.JSON.Error(w, http.StatusBadRequest, "Failed to decode body")
		return
	}
		`, from, bodyStruct)
		return goCode
	} else if singleExists {
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

	return ""
}

type SqlcParams struct {
	structured map[string]map[string]string
	single     map[string]string
}

func parseSqlcFile(filePath string) (SqlcParams, error) {
	params := SqlcParams{}
	structParams := make(map[string]map[string]string)
	singleParams := make(map[string]string)

	var currentStruct string
	var insideStruct bool

	file, err := os.Open(filePath)
	if err != nil {
		return params, err
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
		return params, err
	}

	params.single = singleParams
	params.structured = structParams
	return params, nil
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

func parseInjectionsFile(filePath string) (SqlcParams, error) {
	params := SqlcParams{
		structured: make(map[string]map[string]string),
		single:     make(map[string]string),
	}

	var currentStruct string
	var insideStruct bool

	file, err := os.Open(filePath)
	if err != nil {
		return params, err
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
				params.structured[currentStruct] = make(map[string]string)
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
				params.structured[currentStruct][fieldName] = fieldType
			}
			continue
		}

		// Detect function signature
		if strings.HasPrefix(line, "func (i *Injections)") {
			start := strings.Index(line, "(")
			end := strings.LastIndex(line, ")")
			if start != -1 && end != -1 {
				funcSig := line[start+1 : end]
				paramsList := strings.Split(funcSig, ",")

				// Parse params (skipping ctx)
				for _, param := range paramsList {
					paramParts := strings.Fields(strings.TrimSpace(param))
					if len(paramParts) == 2 {
						paramType := paramParts[1]
						if _, exists := params.structured[paramType]; !exists {
							// Only add if it's a struct type in this file
							params.structured[paramType] = make(map[string]string)
						}
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return params, err
	}

	return params, nil
}
