package generator

import (
	"fmt"
	"master-gen/parser"
	"master-gen/writer"
	"strings"
)

func getNamespace(endpoint string) (string, error) {
	parts := strings.Split(endpoint, "/")
	if len(parts) < 3 {
		return "", fmt.Errorf("invalid endpoint format: %s", endpoint)
	}
	return parts[2], nil
}

func packageAndImports(hasQueries, hasInjections, hasBody bool) string {

	code := `
package handler

import (
	"context"
	"net/http"
	"app/genesis/types"
	`
	if hasBody {
		code += `"encoding/json"`
	}

	if hasQueries {
		code += `
	"time"`
	}

	if hasInjections {
		code += `
	"app/genesis/injection"`

	}

	code += `
)
`

	return code
}

func genHandlers(ops parser.Bliss, path string) error {
	handlers := make(map[string]string)

	hasQuery := false
	hasInjection := false
	hasBody := false
	for _, op := range ops.Operations {
		if op.Body != nil {
			hasBody = true
		}

		op.Query = strings.TrimSpace(op.Query)
		if op.Query != "" {
			hasQuery = true

		}
		op.Handler = strings.TrimSpace(op.Handler)
		if op.Handler != "" {
			hasInjection = true

		}
	}

	for _, op := range ops.Operations {
		namespace, err := getNamespace(op.Endpoint)
		if err != nil {
			return fmt.Errorf("failed to extract namespace: %w", err)
		}
		code := generateHandler(op)

		if existingCode, ok := handlers[namespace]; ok {
			handlers[namespace] = existingCode + code
		} else {
			handlers[namespace] = packageAndImports(hasQuery, hasInjection, hasBody) + code
		}
	}

	for namespace, code := range handlers {
		filePath := fmt.Sprintf(path+"/handler/%s.go", namespace)
		err := writer.WriteFile(filePath, code)
		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}
	}

	return nil
}
