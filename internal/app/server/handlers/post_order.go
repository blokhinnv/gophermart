package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/blokhinnv/gophermart/internal/app/accrual"
	"github.com/blokhinnv/gophermart/internal/app/database"
	"github.com/blokhinnv/gophermart/internal/app/database/ordertracker"
	"github.com/jackc/pgx/v5"
	"golang.org/x/sync/errgroup"
)

type PostOrder struct {
	db                        database.Service
	tracker                   ordertracker.Tracker
	accrualSystem             accrual.Service
	accrualSystemPoolInterval time.Duration
	ctx                       context.Context
	wg                        *sync.WaitGroup
}

type postOrderBody struct {
	OrderID string `valid:"luhn,required"`
}

type accrualSystemResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

const postOrderContentType = "text/plain"

func NewPostOrder(
	db database.Service,
	nWorkers int,
	serverCtx context.Context,
	accrualSystem accrual.Service,
	accrualSystemPoolInterval time.Duration,
) *PostOrder {

	o := PostOrder{
		db:                        db,
		tracker:                   db.Tracker(),
		accrualSystem:             accrualSystem,
		accrualSystemPoolInterval: accrualSystemPoolInterval,
		wg:                        new(sync.WaitGroup),
	}
	g, _ := errgroup.WithContext(serverCtx)
	for i := 0; i < nWorkers; i++ {
		o.wg.Add(1)
		g.Go(o.Loop)
	}
	// теперь это контекст из RunServer
	// он отменяется по сигналу
	o.ctx = serverCtx
	return &o
}

func (h *PostOrder) Loop() error {
	ticker := time.NewTicker(h.accrualSystemPoolInterval)
	for {
		select {
		// если контекст сервера завершен, то завершаем все горутины
		case <-h.ctx.Done():
			// это убьёт всю errgroup
			log.Println("Shutting down Loop goroutine...")
			h.wg.Done()
			return ErrServerShutdown
		case <-ticker.C:
			err := h.LoopIteration()
			if err != nil {
				log.Printf("Error while LoopIteration: %v", err)
				return err
			}
		}
	}
}

func (h *PostOrder) LoopIteration() error {
	// забираем задачу
	task, err := h.tracker.Acquire(h.ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return err
	}
	// делаем запрос к системе расчета баллов
	res, err := h.accrualSystem.GetOrderInfo(task.OrderID)
	if err != nil {
		// если заспамили - возвращаем задачу в работу и отдыхаем несколько секунд
		// оказалось, что при 429 черный ящик возвращает заголовок Retry-After
		// переписал так, чтобы его можно было использовать
		var tmrErr *accrual.ErrTooManyRequests
		if errors.As(err, &tmrErr) {
			err = h.tracker.UpdateStatusAndRelease(h.ctx, task.StatusID, task.OrderID)
			if err != nil {
				return err
			}
			time.Sleep(tmrErr.RetryAfter)
			return nil
		}
		return err
	}

	resp := accrualSystemResponse{}
	json.Unmarshal(res, &resp)
	// обновить запись о заказе
	err = h.db.UpdateOrderStatus(context.Background(), task.OrderID, resp.Status)
	if err != nil {
		return err
	}
	// проверить готовность
	switch {
	case resp.Status == "PROCESSED":
		// если заказ обработан - добавим запись с баллами и удалим из очереди на обработку
		err = h.db.AddAccrualRecord(context.Background(), task.OrderID, resp.Accrual)
		if err != nil {
			return err
		}
		err = h.tracker.Delete(h.ctx, task.OrderID)
		if err != nil {
			return err
		}
	case resp.Status == "REGISTERED" || resp.Status == "PROCESSING":
		// если не обработан - возвращаем в работу
		err = h.tracker.UpdateStatusAndRelease(
			h.ctx,
			database.STATUSES[resp.Status],
			task.OrderID,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *PostOrder) WaitDone() {
	h.wg.Wait()
}

func (h *PostOrder) ReadBody(r *http.Request) (*postOrderBody, int, error) {
	bodyReader := func(bodyBytes []byte) (any, error) {
		body := postOrderBody{OrderID: string(bodyBytes)}
		return &body, nil
	}
	body, err := ReadBodyWithBodyReader(r, postOrderContentType, bodyReader)
	if err != nil {
		if errors.Is(err, ErrNotValid) {
			return nil, http.StatusUnprocessableEntity, err
		}
		return nil, http.StatusBadRequest, err
	}
	if bodyTyped, ok := body.(*postOrderBody); ok {
		return bodyTyped, http.StatusOK, nil
	} else {
		return nil, http.StatusInternalServerError, nil
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
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = h.db.AddOrder(ctx, body.OrderID, userID)
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
	h.tracker.Add(ctx, body.OrderID)
	w.WriteHeader(http.StatusAccepted)
}
