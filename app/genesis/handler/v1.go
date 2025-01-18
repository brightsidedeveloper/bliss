package handler

import (
	"app/genesis/injection"
	"app/genesis/types"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

func (h *Handler) GetAha(w http.ResponseWriter, r *http.Request) {
	queryParams := types.AhaQuery{
		Anything: r.URL.Query().Get("anything"),
		Example:  r.URL.Query().Get("example"),
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

func (h *Handler) PostAha3(w http.ResponseWriter, r *http.Request) {
	body := types.Aha3Body{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.JSON.Error(w, http.StatusBadRequest, "Failed to decode body")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT count FROM public.testies t WHERE t.size = $1"

	rows, err := h.DB.QueryContext(ctx, query)
	if err != nil {
		h.JSON.Error(w, http.StatusInternalServerError, "Failed to query user")
		return
	}
	defer rows.Close()

	res := make([]types.Aha3Row, 0)
	for rows.Next() {
		row := types.Aha3Row{}
		err = rows.Scan()
		if err != nil {
			h.JSON.Error(w, http.StatusInternalServerError, "Failed to query user")
			return
		}
		res = append(res, row)
	}

	h.JSON.Success(w, res)
}
