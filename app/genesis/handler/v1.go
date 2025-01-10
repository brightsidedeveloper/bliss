package handler

import (
	"context"
	"net/http"
	"time"
)

func (h *Handler) GetAha(w http.ResponseWriter, r *http.Request) {

	type AhaQuery struct {
		Name string `json:"name"`
	}

	queryParams := AhaQuery{
		Name: r.URL.Query().Get("name"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT name, msg FROM public.profile p WHERE p.name = $1"

	type AhaRow struct {
		Name string `json:"name"`
		Msg  string `json:"msg"`
	}

	res := AhaRow{}
	err := h.DB.QueryRowContext(ctx, query, queryParams.Name).Scan(&res.Name, &res.Msg)
	if err != nil {
		h.JSON.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.JSON.Success(w, res)
}
