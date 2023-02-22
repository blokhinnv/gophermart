package accrual

import (
	"fmt"
	"time"
)

// var ErrTooManyRequests = errors.New("too many requests for the accrual system")

type ErrTooManyRequests struct {
	RetryAfter time.Duration
}

func (e *ErrTooManyRequests) Error() string {
	return fmt.Sprintf("too many requests for the accrual system; wait %v", e.RetryAfter)
}

func NewErrTooManyRequests(retryAfter int) error {
	return &ErrTooManyRequests{
		RetryAfter: time.Duration(retryAfter) * time.Second,
	}
}
