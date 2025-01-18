
package handler

import (
	"context"
	"net/http"
	"solar-system/genesis/types"
	"encoding/json"
	"time"
)

func (h *Handler) GetAha(w http.ResponseWriter, r *http.Request) {
	queryParams := types.AhaQuery{
		Example: r.URL.Query().Get("example"),
		Anything: r.URL.Query().Get("anything"),
	}
	
				
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT name, email FROM public.profile p WHERE p.id = $1 AND p.name = $2"

	
	res := types.AhaRow{}
	err := h.DB.QueryRowContext(ctx, query, queryParams.Example, queryParams.Anything).Scan(&res.Name, &res.Email)
	if err != nil {
		h.JSON.Error(w, http.StatusInternalServerError, "Failed to query user")
		return
	}

	h.JSON.Success(w, res)
}	

func (h *Handler) DeleteAha3(w http.ResponseWriter, r *http.Request) {
	body := types.Aha3Body{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.JSON.Error(w, http.StatusBadRequest, "Failed to decode body")
		return
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "DELETE FROM public.testies t WHERE t.size = $1"

	
	res := types.Aha3Row{}
	err := h.DB.QueryRowContext(ctx, query, body.Size)
	if err != nil {
		h.JSON.Error(w, http.StatusInternalServerError, "Failed to query user")
		return
	}

	h.JSON.Success(w, res)
}	
