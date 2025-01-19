package generator

import (
	"fmt"
	"master-gen/internal/parser"
	"master-gen/internal/writer"
	"path"
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

	code := `
package handler

import (
	"context"
	"net/http"
	"encoding/json"
	`

	if hasQueries {
		code += `
	"time"
	"bliss-server/genesis/queries"`
	}

	if hasInjections {
		code += `
	"bliss-server/genesis/injection"`

	}

	code += `
)
`

	return code
}

func genHandlers(g *Generator, ops parser.Bliss, dest string) error {
	handlers := make(map[string]string)
	structParams, singleParams, err := parseSqlcFile(path.Join(g.ServerPath, "genesis/queries/queries.sql.go"))
	if err != nil {
		return err
	}
	hasQuery := false
	hasInjection := false
	for _, op := range ops.Operations {

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
		code := generateHandler(op, structParams, singleParams)

		if existingCode, ok := handlers[namespace]; ok {
			handlers[namespace] = existingCode + code
		} else {
			handlers[namespace] = packageAndImports(hasQuery, hasInjection) + code
		}
	}

	for namespace, code := range handlers {
		filePath := fmt.Sprintf(dest+"/handler/%s.go", namespace)
		err := writer.WriteFile(filePath, code)
		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}
	}

	return nil
}
