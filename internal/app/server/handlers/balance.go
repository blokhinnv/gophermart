package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/blokhinnv/gophermart/internal/app/database"
)

type Balance struct {
	db database.Service
}

func (h *Balance) Handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	balance, err := h.db.GetBalance(ctx, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	balanceEncoded, err := json.Marshal(balance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(balanceEncoded)
}
