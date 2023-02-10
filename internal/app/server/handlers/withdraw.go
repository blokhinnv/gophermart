package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/blokhinnv/gophermart/internal/app/database"
	"github.com/go-chi/jwtauth/v5"
)

type Withdraw struct {
	db database.Service
}

type withdrawRequestBody struct {
	OrderID string  `valid:"orderID,required" json:"order"`
	Sum     float64 `valid:",required"        json:"sum"`
}

const withdrawContentType = "application/json"

func (h *Withdraw) ReadBody(r *http.Request) (*withdrawRequestBody, int, error) {
	bodyReader := func(bodyBytes []byte) (any, error) {
		body := withdrawRequestBody{}
		if err := json.Unmarshal(bodyBytes, &body); err != nil {
			return nil, fmt.Errorf(
				"%w: incorrent body (error while unmarshaling)",
				ErrIncorrectRequest,
			)
		}
		return &body, nil
	}
	body, err := ReadBodyWithBodyReader(r, withdrawContentType, bodyReader)
	if err != nil {
		if errors.Is(err, ErrNotValid) {
			return nil, http.StatusUnprocessableEntity, err
		}
		return nil, http.StatusBadRequest, err
	}
	if bodyTyped, ok := body.(*withdrawRequestBody); ok {
		return bodyTyped, http.StatusOK, nil
	} else {
		return nil, http.StatusInternalServerError, nil
	}
}

func (h *Withdraw) Handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// читаем тело
	body, status, err := h.ReadBody(r)
	// если запрос некорретный - заканчиваем работу
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), status)
		return
	}
	_, claims, _ := jwtauth.FromContext(ctx)
	userID, ok := claims["user_id"].(float64)
	if !ok {
		http.Error(w, "no int user_id in claims", http.StatusInternalServerError)
		return
	}
	balance, err := h.db.GetBalance(ctx, int(userID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if balance.Current.Float64 < body.Sum {
		http.Error(
			w,
			fmt.Errorf("%w: userID=%v request=%+v", ErrNotEnoughBalance, userID, body).Error(),
			http.StatusPaymentRequired,
		)
		return
	}
	err = h.db.AddWithdrawalRecord(ctx, body.OrderID, body.Sum)
	if err != nil {
		if errors.Is(err, database.ErrMissingOrderID) {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
