package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/blokhinnv/gophermart/internal/app/database"
)

type GetOrder struct {
	db database.Service
}

func (h *GetOrder) Handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	orders, err := h.db.FindOrdersByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, database.ErrEmptyResult) {
			http.Error(w, err.Error(), http.StatusNoContent)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ordersEncoded, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(ordersEncoded)
}
