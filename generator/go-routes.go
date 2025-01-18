package generator

import (
	"fmt"
	"master-gen/parser"
	"master-gen/writer"
)

func MountRoutesFileContent(ops parser.AhaJSON) error {

	operationNames := make(map[string]int)
	for _, op := range ops.Operations {
		operationNames[op.Name]++
		if operationNames[op.Name] > 1 {
			return fmt.Errorf("operation name %s is duplicated", op.Name)
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
			return fmt.Errorf("endpoint %s with method %s is duplicated", op.Endpoint, op.Method)
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

	writer.WriteFile("./app/genesis/routes/routes.go", goCode)

	return nil
}
