package handler

import (
	"app/genesis/injection"
	"app/genesis/types"
	"context"
	"net/http"
	"time"
)

func (h *Handler) GetAha(w http.ResponseWriter, r *http.Request) {

	queryParams := types.AhaQuery{
		Example: r.URL.Query().Get("example"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT name, email FROM public.profile p WHERE p.id = $1 AND p.name = $2"

	res := types.AhaRow{}
	err := h.DB.QueryRowContext(ctx, query, queryParams.Example).Scan(&res.Email, &res.Name)
	if err != nil {
		h.JSON.Error(w, http.StatusInternalServerError, "Failed to query user")
		return
	}

	h.JSON.Success(w, res)
}

func (h *Handler) GetAha2(w http.ResponseWriter, r *http.Request) {

	queryParams := types.Aha2Query{
		Example: r.URL.Query().Get("example"),
	}

	res, err := injection.CheckCoolStatus(w, r, queryParams)
	if err != nil {
		h.JSON.Error(w, http.StatusInternalServerError, "Failed to query user")
		return
	}

	h.JSON.Success(w, res)
}
