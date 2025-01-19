package generator

import (
	"fmt"
	"master-gen/internal/parser"
	"master-gen/internal/writer"
)

func genMountRoutes(ops parser.Bliss, path string) error {
	operationNames := make(map[string]int)
	for _, op := range ops.Operations {
		operationNames[op.Query]++
		if operationNames[op.Query] > 1 {
			return fmt.Errorf("operation query %s is duplicated", op.Query)
		}
	}

	endpointPaths := make(map[string]int)
	for _, op := range ops.Operations {
		endpointPaths[op.Endpoint]++
		if endpointPaths[op.Endpoint] > 1 {
			return fmt.Errorf("endpoint %s is duplicated", op.Endpoint)
		}
	}

	goCode := `package routes
	
import (
	"bliss-server/genesis/handler"
	"github.com/go-chi/chi/v5"
)

func MountRoutes(r *chi.Mux, h *handler.Handler) {
`
	for _, op := range ops.Operations {
		goCode += `
	r.Post("` + op.Endpoint + `", h.` + op.Query + `)`
	}
	goCode += `
}
`

	// Write the generated code to the specified file
	if err := writer.WriteFile(path+"/routes/routes.go", goCode); err != nil {
		return fmt.Errorf("failed to write routes file: %w", err)
	}

	return nil
}
