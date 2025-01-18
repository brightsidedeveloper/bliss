
package handler

import (
	"app/genesis/types"
	"context"
	"net/http"
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

}

func (h *Handler) GetSuperTest(w http.ResponseWriter, r *http.Request) {
	queryParams := types.SuperTestQuery{
		Example: r.URL.Query().Get("example"),
	}
	
				
	res, err := injection.CheckCoolStatus(w, r, queryParams)
	if err != nil {
		h.JSON.Error(w, http.StatusInternalServerError, "Failed to query user")
		return
	}
	
	h.JSON.Success(w, res)
}	

func (h *Handler) GetAha3(w http.ResponseWriter, r *http.Request) {
	queryParams := types.Aha3Query{
		Size: r.URL.Query().Get("size"),
	}
	
				
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT count FROM public.testies t WHERE t.size = $1"

	rows, err := h.DB.QueryContext(ctx, query, queryParams.Size)
	if err != nil {
		h.JSON.Error(w, http.StatusInternalServerError, "Failed to query user")
		return
	}
	defer rows.Close()

	res := make([]types.Aha3Row, 0)
	for rows.Next() {
		row := types.Aha3Row{}
		err = rows.Scan(queryParams.Size)
		if err != nil {
			h.JSON.Error(w, http.StatusInternalServerError, "Failed to query user")
			return
		}
		res = append(res, row)
	}

}
