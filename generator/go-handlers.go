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

func packageAndImports(hasQueries, hasInjections bool) string {

	if hasInjections && hasQueries {
		return `
package handler

import (
	"app/genesis/injection"
	"app/genesis/types"
	"context"
	"net/http"
	"time"
)
	`
	}
	if hasInjections {
		return `
package handler

import (
	"app/genesis/injection"
	"app/genesis/types"
	"context"
	"net/http"
)
`
	}
	if hasQueries {
		return `
package handler

import (
	"app/genesis/types"
	"context"
	"net/http"
	"time"
)
`
	}
	return ""
}

func Handlers(ops parser.AhaJSON) error {
	handlers := make(map[string]string)

	for _, op := range ops.Operations {
		namespace, err := getNamespace(op.Endpoint)
		if err != nil {
			return fmt.Errorf("failed to extract namespace: %w", err)
		}
		var code string
		switch op.Method {
		case "Get":
			code = generateGetHandler(op)
		default:
			continue
		}
		if existingCode, ok := handlers[namespace]; ok {
			handlers[namespace] = existingCode + code
		} else {
			handlers[namespace] = packageAndImports(op.QueryParams != nil, op.Handler != "") + code
		}
	}

	for namespace, code := range handlers {
		filePath := fmt.Sprintf("./app/genesis/handler/%s.go", namespace)
		err := writer.WriteFile(filePath, code)
		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}
	}

	return nil
}
