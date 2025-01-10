package handler

import (
	"context"
	"net/http"
	"encoding/json"
	"time"
)

	
func (h *Handler) GetAha(w http.ResponseWriter, r *http.Request) {
		
	type AhaQuery struct {
		Example string `json:"example"`
	}

	var queryParams AhaQuery
			
	err := json.NewDecoder(r.Body).Decode(&queryParams)
	if err != nil {
		h.JSON.ValidationError(w, "Bad request")
		return
	}
				
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT name, email FROM public.profile p WHERE p.id = $1 AND p.name = $2"

	type AhaRow struct {
		Email string `json:"email"`
		Name string `json:"name"`
	}
	
	res := AhaRow{}
	err = h.DB.QueryRowContext(ctx, query, queryParams.Example).Scan(&res.Name, &res.Email)
	if err != nil {
		h.JSON.Error(w, http.StatusInternalServerError, "Failed to query user")
		return
	}

	h.JSON.Success(w, res)
}