package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/blokhinnv/gophermart/internal/app/database"
	"github.com/go-chi/jwtauth/v5"
)

type GetOrder struct {
	db database.Service
}

func (h *GetOrder) Handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, claims, _ := jwtauth.FromContext(ctx)
	userID, ok := claims["user_id"].(float64)
	if !ok {
		http.Error(w, "no int user_id in claims", http.StatusInternalServerError)
		return
	}
	orders, err := h.db.FindOrdersByUserID(ctx, int(userID))
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
