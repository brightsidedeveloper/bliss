package generator

import (
	"master-gen/internal/parser"
	"master-gen/internal/writer"
)

func genGoTypes(ops parser.Bliss, path string) error {

	goCode := `package types

	`

	writer.WriteFile(path+"/types/types.go", goCode)

	return nil
}
