package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/blokhinnv/gophermart/internal/app/accrual"
	"github.com/blokhinnv/gophermart/internal/app/database"
	"github.com/go-chi/jwtauth/v5"
)

type PostOrder struct {
	db            database.Service
	c             chan string
	accrualSystem accrual.Service
}

type postOrderBody struct {
	OrderID string `valid:"orderID,required"`
}

type accrualSystemResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

const orderContentType = "text/plain"

func NewPostOrder(
	db database.Service,
	cSize, nWorkers int,
	accrualSystem accrual.Service,
) *PostOrder {
	o := PostOrder{
		db:            db,
		c:             make(chan string, cSize),
		accrualSystem: accrualSystem,
	}
	for i := 0; i < nWorkers; i++ {
		go o.Loop()
	}
	return &o
}

func (h *PostOrder) ReadBody(r *http.Request) (*postOrderBody, int, error) {
	// проверим content-type
	if r.Header.Get("Content-Type") != orderContentType {
		return nil, http.StatusBadRequest, fmt.Errorf(
			"%w: incorrect content type",
			ErrIncorrectRequest,
		)
	}
	// проверим содержимое
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil || len(bodyBytes) == 0 {
		return nil, http.StatusBadRequest, fmt.Errorf(
			"%w: incorrent body (error while reading)",
			ErrIncorrectRequest,
		)
	}
	body := postOrderBody{OrderID: string(bodyBytes)}
	if validated, err := govalidator.ValidateStruct(body); err != nil || !validated {
		return nil, http.StatusUnprocessableEntity, fmt.Errorf(
			"%w: incorrent body (error while validating) %v",
			ErrIncorrectRequest,
			err.Error(),
		)
	}
	return &body, http.StatusOK, nil

}

func (h *PostOrder) Loop() {
	for {
		orderID := <-h.c
		// делаем запрос к системе расчета баллов
		res, err := h.accrualSystem.GetOrderInfo(orderID)
		if err != nil {
			log.Printf("Error while processing order %v: %v\n", orderID, err.Error())
			continue
		}
		// TODO: remove after debug
		time.Sleep(10 * time.Second)
		resp := accrualSystemResponse{}
		json.Unmarshal(res, &resp)
		// обновить запись о заказе
		h.db.UpdateOrderStatus(context.Background(), orderID, resp.Status)
		// проверить готовность
		switch {
		case resp.Status == "PROCESSED":
			h.db.AddAccrualRecord(context.Background(), orderID, resp.Accrual)
		case resp.Status == "REGISTERED" || resp.Status == "PROCESSING":
			h.c <- orderID
		}
	}
}

func (h *PostOrder) Handler(w http.ResponseWriter, r *http.Request) {
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
	err = h.db.AddOrder(ctx, body.OrderID, int(userID))
	if err != nil {
		switch {
		case errors.Is(err, database.ErrOrderAlreadyAddedByThisUser):
			http.Error(w, err.Error(), http.StatusOK)
		case errors.Is(err, database.ErrOrderAlreadyAddedByOtherUser):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	h.c <- body.OrderID
	w.WriteHeader(http.StatusAccepted)
}
