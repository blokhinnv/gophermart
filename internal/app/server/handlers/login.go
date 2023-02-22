package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/blokhinnv/gophermart/internal/app/auth"
	"github.com/blokhinnv/gophermart/internal/app/database"
)

type Login struct {
	LogReg
}

func (h *Login) Handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// читаем тело
	body, status, err := h.ReadBody(r)
	// если запрос некорретный - заканчиваем работу
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	user, err := h.db.FindUser(ctx, body.Login, body.Password)
	if err != nil {
		// если не нашли пользователя - не можем авторизоваться
		if errors.Is(err, database.ErrUserNotFound) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// если хэш предоставленного пароля не совпадает с тем, что лежит в БД -
	// не можем авторизоваться
	if hash := auth.GenerateHash(body.Password, user.Salt); hash != user.HashedPassword {
		err := fmt.Errorf("%v: %v %v", ErrIncorrectCredentials, body.Login, body.Password)
		http.Error(
			w,
			err.Error(),
			http.StatusUnauthorized,
		)
		return
	}
	token := auth.GenerateJWTToken(user, h.signingKey, h.expireDuration)
	tokenSign, err := token.SignedString(h.signingKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", fmt.Sprintf("Bearer: %v", tokenSign))
	w.WriteHeader(http.StatusOK)
}
