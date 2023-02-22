package accrual

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
)

type AccrualService struct {
	client *resty.Client
}

func NewAccrualService(accrualSystemAddr string) *AccrualService {
	client := resty.New()
	client.SetBaseURL(fmt.Sprintf("%v/api/orders/", accrualSystemAddr))
	return &AccrualService{client: client}
}

func (s *AccrualService) GetOrderInfo(orderID string) ([]byte, error) {
	res, err := s.client.R().Get(orderID)
	if err != nil {
		return nil, err
	}

	if res.StatusCode() == http.StatusTooManyRequests {
		retryAfter, err := strconv.Atoi(res.Header().Get("Retry-After"))
		if err != nil {
			return nil, err
		}
		return nil, NewErrTooManyRequests(retryAfter)
	}
	return res.Body(), nil
}
