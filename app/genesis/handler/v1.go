
package handler

import (
	"context"
	"net/http"
	"solar-system/genesis/types"
	
	"time"
)

func (h *Handler) GetExample(w http.ResponseWriter, r *http.Request) {
	queryParams := types.ExampleQuery{
		Example: r.URL.Query().Get("example"),
	}
	
				
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT name FROM public.example e WHERE e.example = $1"

	
	res := types.ExampleRow{}
	err := h.DB.QueryRowContext(ctx, query, queryParams.Example).Scan(&res.Name)
	if err != nil {
		h.JSON.Error(w, http.StatusInternalServerError, "Failed to query user")
		return
	}

	h.JSON.Success(w, res)
}	
