package injection

import (
	"app/genesis/types"
	"net/http"
)

// Add Injection Context
func CheckCoolStatus(w http.ResponseWriter, r *http.Request, queryParams types.Aha2Query) (types.Aha2Res, error) {

	//I Can do anything here
	return types.Aha2Res{Cool: true}, nil
}
