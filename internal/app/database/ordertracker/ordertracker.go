package ordertracker

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBTracker struct {
	conn *pgxpool.Pool
}

func NewDBTracker(conn *pgxpool.Pool) *DBTracker {
	return &DBTracker{conn: conn}
}

func (q *DBTracker) Acquire(ctx context.Context) (*Task, error) {
	var orderID string
	var statusID int
	err := q.conn.QueryRow(ctx, acquireSQL).Scan(&orderID, &statusID)
	if err != nil {
		return nil, err
	}
	return &Task{OrderID: orderID, StatusID: statusID}, nil
}

func (q *DBTracker) UpdateStatusAndRelease(
	ctx context.Context,
	newStatusID int,
	orderID string,
) error {
	_, err := q.conn.Exec(ctx, updateAndReleaseSQL, newStatusID, orderID)
	return err
}

func (q *DBTracker) Add(ctx context.Context, orderID string) error {
	_, err := q.conn.Exec(ctx, addSQL, orderID)
	return err
}

func (q *DBTracker) Delete(ctx context.Context, orderID string) error {
	_, err := q.conn.Exec(ctx, deleteSQL, orderID)
	return err
}
