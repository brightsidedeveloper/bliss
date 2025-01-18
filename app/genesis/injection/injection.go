package injection

import (
	"app/genesis/types"
	"net/http"
)

// Add Injection Context
func CheckCoolStatus(w http.ResponseWriter, r *http.Request, queryParams types.SuperTestQuery) (types.SuperTestRes, error) {

	//I Can do anything here
	return types.SuperTestRes{Cool: true}, nil
}
