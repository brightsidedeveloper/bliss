package generator

import (
	"fmt"
	"master-gen/internal/writer"
)

func genTsTypes(params ParamStrings, responses ResponseStrings, path string) error {

	code := ``

	for _, p := range params {
		code += fmt.Sprintf(`
		%s`, p.interfaceType)
	}

	for _, r := range responses {
		code += fmt.Sprintf(`
		%s`, r.interfaceType)
	}

	return writer.WriteFile(path+"/types.ts", code)
}
