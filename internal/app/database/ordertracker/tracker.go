package ordertracker

import "context"

type Task struct {
	OrderID  string
	StatusID int
}

type Tracker interface {
	Acquire(ctx context.Context) (*Task, error)
	UpdateStatusAndRelease(
		ctx context.Context,
		newStatusID int,
		orderID string,
	) error
	Add(ctx context.Context, orderID string) error
	Delete(ctx context.Context, orderID string) error
}
