package generator

import (
	"master-gen/parser"
	"master-gen/util"
	"strings"
)

func generateGetHandler(op parser.Operation) string {
	switch op.ResType {
	case "row":
		return generateGetQueryRowHandler(op)
	case "rows":
		return generateGetQueryRowsHandler(op)
	case "custom":
		return generateGetCustomHandler(op)
	default:
		return ""
	}
}

func generateQueryParams(name string, params parser.QueryParams) string {
	if params == nil {
		return ""
	}

	goCode := `
	queryParams := types.` + name + `Query{`
	for k := range params {
		goCode += `
		` + util.Capitalize(k) + `: r.URL.Query().Get("` + k + `"),`
	}
	goCode += `
	}
	
				`
	return goCode
}

func generateGetQueryRow(name, queryStr string, params parser.QueryParams, res parser.Response) string {

	inserts := make([]string, 0)
	for k := range params {
		inserts = append(inserts, "queryParams."+util.Capitalize(k))
	}

	query, _ := processQuery(queryStr)

	goCode := `
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "` + query + `"

	`

	scan := make([]string, 0)
	for k := range res {
		scan = append(scan, "&res."+util.Capitalize(k))
	}
	goCode += `
	res := types.` + name + `Row{}
	err := h.DB.QueryRowContext(ctx, query, ` + strings.Join(inserts, ", ") + `).Scan(` + strings.Join(scan, ", ") + `)
	if err != nil {
		h.JSON.Error(w, http.StatusInternalServerError, "Failed to query user")
		return
	}
`

	return goCode
}

func generateGetQueryRowHandler(op parser.Operation) string {
	goCode := ""
	goCode += `
func (h *Handler) ` + op.Method + op.Name + `(w http.ResponseWriter, r *http.Request) {`
	if op.QueryParams != nil {
		goCode += generateQueryParams(op.Name, op.QueryParams)
	}
	goCode += generateGetQueryRow(op.Name, op.Query, op.QueryParams, op.Res)

	goCode += `
}
`
	return goCode
}

func generateGetQueryRows(name, queryStr string, params parser.QueryParams) string {

	inserts := make([]string, 0)
	for k := range params {
		inserts = append(inserts, "queryParams."+util.Capitalize(k))
	}

	query, _ := processQuery(queryStr)

	goCode := `
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "` + query + `"

	rows, err := h.DB.QueryContext(ctx, query, ` + strings.Join(inserts, ", ") + `)
	if err != nil {
		h.JSON.Error(w, http.StatusInternalServerError, "Failed to query user")
		return
	}
	defer rows.Close()

	res := make([]types.` + name + `Row, 0)
	for rows.Next() {
		row := types.` + name + `Row{}
		err = rows.Scan(` + strings.Join(inserts, ", ") + `)
		if err != nil {
			h.JSON.Error(w, http.StatusInternalServerError, "Failed to query user")
			return
		}
		res = append(res, row)
	}
`

	return goCode
}

func generateGetQueryRowsHandler(op parser.Operation) string {
	goCode := ""
	goCode += `
func (h *Handler) ` + op.Method + op.Name + `(w http.ResponseWriter, r *http.Request) {`
	if op.QueryParams != nil {
		goCode += generateQueryParams(op.Name, op.QueryParams)
	}
	goCode += generateGetQueryRows(op.Name, op.Query, op.QueryParams)

	goCode += `
}
`
	return goCode
}

func generateGetCustomHandler(op parser.Operation) string {
	goCode := ""
	goCode += `
func (h *Handler) ` + op.Method + op.Name + `(w http.ResponseWriter, r *http.Request) {`
	if op.QueryParams != nil {
		goCode += generateQueryParams(op.Name, op.QueryParams)
	}
	goCode += `
	res, err := injection.` + processHandler(op.Handler, op.QueryParams) + `
	if err != nil {
		h.JSON.Error(w, http.StatusInternalServerError, "Failed to query user")
		return
	}
	`
	goCode += `
	h.JSON.Success(w, res)
}	
`
	return goCode
}
