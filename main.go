package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	file, err := os.Open("./app/aha.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var ops AhaJSON
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&ops)
	if err != nil {
		log.Fatal(err)
	}

	routesFileContent, err := generateRoutesFileContent(ops)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile("./app/genesis/routes/routes.go", []byte(routesFileContent), 0644)
	if err != nil {
		log.Fatal(err)
	}
	handlerFilesContent, err := generateHandlerFilesContent(ops)
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range handlerFilesContent {
		err = os.WriteFile("./app/genesis/handler/"+k+".go", []byte(v), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
	typesFileContent, err := generateTypesFileContent(ops)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile("./app/genesis/types/types.go", []byte(typesFileContent), 0644)
	if err != nil {
		log.Fatal(err)
	}

}

func generateRoutesFileContent(ops AhaJSON) (string, error) {

	operationNames := make(map[string]int)
	for _, op := range ops.Operations {
		operationNames[op.Name]++
		if operationNames[op.Name] > 1 {
			return "", fmt.Errorf("operation name %s is duplicated", op.Name)
		}
	}
	endpointMethods := make(map[string]map[string]int)
	for _, op := range ops.Operations {
		_, ok := endpointMethods[op.Method]
		if ok {
			endpointMethods[op.Method][op.Endpoint]++
		} else {
			endpointMethods[op.Method] = make(map[string]int)
			endpointMethods[op.Method][op.Endpoint]++
		}
		if endpointMethods[op.Method][op.Endpoint] > 1 {
			return "", fmt.Errorf("endpoint %s with method %s is duplicated", op.Endpoint, op.Method)
		}
	}

	goCode := `package routes
	
	import (
	"app/genesis/handler"

	"github.com/go-chi/chi/v5"
)

func MountRoutes(r *chi.Mux, h *handler.Handler) {
	`
	for _, op := range ops.Operations {
		goCode += `
		r.` + op.Method + `("` + op.Endpoint + `", h.` + op.Method + op.Name + `)`
	}
	goCode += `

}
	`

	return goCode, nil
}

func capitalize(s string) string {
	return strings.ToUpper(string(s[0])) + s[1:]
}

func generateTypesFileContent(ops AhaJSON) (string, error) {

	goCode := `package types

	`
	for _, op := range ops.Operations {
		theType := "Row"
		if op.Handler != "" {
			theType = "Res"
		}
		goCode += `
type ` + op.Name + theType + ` struct {`
		for k, t := range op.Res {

			goCode += `
	` + capitalize(k) + ` ` + t + " `" + `json:"` + k + `"` + "`"
		}
		goCode += `
}
	`
		if op.QueryParams != nil {
			goCode += `
type ` + op.Name + `Query struct {`
			for k, t := range op.QueryParams {
				goCode += `
	` + capitalize(k) + ` ` + t + " `" + `json:"` + k + `"` + "`"
			}
			goCode += `
}

	`
		}
	}

	return goCode, nil
}

func generateHandlerFilesContent(ops AhaJSON) (map[string]string, error) {

	initCode := `package handler

import (
	"app/genesis/injection"
	"app/genesis/types"
	"context"
	"net/http"
	"time"
)

	`

	goCodeFiles := make(map[string]string)
	for _, op := range ops.Operations {
		parts := strings.Split(op.Endpoint, "/")
		if len(parts) < 3 {
			return map[string]string{}, fmt.Errorf("invalid endpoint format: %s", op.Endpoint)
		}
		namespace := parts[2]
		goCode := ``

		goCode += `
func (h *Handler) ` + op.Method + op.Name + `(w http.ResponseWriter, r *http.Request) {
		`
		if op.QueryParams != nil {

			goCode += `
	queryParams := types.` + op.Name + `Query{`
			for k := range op.QueryParams {
				goCode += `
		` + capitalize(k) + `: r.URL.Query().Get("` + k + `"),`
			}
			goCode += `
	}
	
				`

		}

		if op.Query != "" {

			inserts := make([]string, 0)
			if op.QueryParams != nil {
				for k := range op.QueryParams {
					inserts = append(inserts, "queryParams."+capitalize(k))
				}
			}
			query, _ := processQuery(op.Query)

			goCode += `
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "` + query + `"

	`

			if op.ResType == "array" {
				scan := make([]string, 0)
				for k := range op.Res {
					scan = append(scan, "&row."+capitalize(k))
				}
				goCode += `
				
	res := make([]types.` + op.Name + `Row, 0)
	rows, err := h.DB.QueryContext(ctx, query, ` + strings.Join(inserts, ", ") + `)
	if err != nil {
		h.JSON.Error(w, http.StatusInternalServerError, "Failed to query users")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var row types.` + op.Name + `Row
		if err := rows.Scan(` + strings.Join(scan, ", ") + `); err != nil {
			h.JSON.Error(w, http.StatusInternalServerError, "Failed to scan users")
			return
		}
		res = append(res, row)
	}
				`
			} else {
				scan := make([]string, 0)
				for k := range op.Res {
					scan = append(scan, "&res."+capitalize(k))
				}
				goCode += `
	res := types.` + op.Name + `Row{}
	err := h.DB.QueryRowContext(ctx, query, ` + strings.Join(inserts, ", ") + `).Scan(` + strings.Join(scan, ", ") + `)
	if err != nil {
		h.JSON.Error(w, http.StatusInternalServerError, "Failed to query user")
		return
	}
`
			}
		} else if op.Handler != "" {
			goCode += `
	res, err := injection.` + processHandler(op.Handler, op.QueryParams) + `
	if err != nil {
		h.JSON.Error(w, http.StatusInternalServerError, "Failed to query user")
		return
	}
`

		} else {
			return map[string]string{}, fmt.Errorf("handler not implemented yet")
		}
		goCode += `
	h.JSON.Success(w, res)
}
`
		code, ok := goCodeFiles[namespace]
		if ok {
			goCodeFiles[namespace] = code + goCode
		} else {
			goCodeFiles[namespace] = initCode + goCode
		}
	}

	return goCodeFiles, nil
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

type AhaJSON struct {
	Operations []Operation `json:"operations"`
}

type QueryParams map[string]string
type Response map[string]string

type Operation struct {
	Name        string      `json:"name"`
	Endpoint    string      `json:"endpoint"`
	Method      string      `json:"method"`
	QueryParams QueryParams `json:"queryParams"`
	Query       string      `json:"query"`
	Handler     string      `json:"handler"`
	ResType     string      `json:"responseType"`
	Res         Response    `json:"response"`
}
