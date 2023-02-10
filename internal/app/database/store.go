package database

import (
	"context"

	"github.com/blokhinnv/gophermart/internal/app/database/ordertracker"
	"github.com/blokhinnv/gophermart/internal/app/models"
)

type Service interface {
	AddUser(ctx context.Context, username, pwd string) (*models.User, error)
	FindUser(ctx context.Context, username, pwd string) (*models.User, error)
	FindOrderByID(ctx context.Context, orderID string) (*models.Order, error)
	AddOrder(ctx context.Context, orderID string, userID int) error
	UpdateOrderStatus(ctx context.Context, orderID, newStatus string) error
	AddAccrualRecord(ctx context.Context, orderID string, sum float64) error
	FindOrdersByUserID(ctx context.Context, userID int) ([]models.Order, error)
	GetBalance(ctx context.Context, userID int) (*models.Balance, error)
	AddWithdrawalRecord(ctx context.Context, orderID string, sum float64, userID int) error
	GetWithdrawals(ctx context.Context, userID int) ([]models.Withdrawal, error)
	Tracker() ordertracker.Tracker
}
