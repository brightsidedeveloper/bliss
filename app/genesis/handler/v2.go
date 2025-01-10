package handler

import (
	"context"
	"net/http"
	"encoding/json"
	"time"
)

	
func (h *Handler) GetAha2(w http.ResponseWriter, r *http.Request) {
		
	type Aha2Query struct {
		Example string `json:"example"`
	}

	var queryParams Aha2Query
			
	err := json.NewDecoder(r.Body).Decode(&queryParams)
	if err != nil {
		h.JSON.ValidationError(w, "Bad request")
		return
	}
				
	h.JSON.Success(w, res)
}