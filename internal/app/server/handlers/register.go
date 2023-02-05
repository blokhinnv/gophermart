package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/blokhinnv/gophermart/internal/app/auth"
	"github.com/blokhinnv/gophermart/internal/app/database"
)

type Register struct {
	LogReg
}

func (h *Register) Handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// читаем тело
	body, status, err := h.ReadBody(r)
	// если запрос некорретный - заканчиваем работу
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}
	err = h.db.AddUser(ctx, body.Login, body.Password)
	if err != nil {
		if errors.Is(err, database.ErrUserAlreadyExists) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token, err := auth.GenerateJWTToken(body.Login, h.signingKey, h.expireDuration)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", fmt.Sprintf("Bearer: %v", token))
	w.WriteHeader(http.StatusOK)
}
