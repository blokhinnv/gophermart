package database

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/blokhinnv/gophermart/internal/app/auth"
	"github.com/blokhinnv/gophermart/internal/app/database/ordertracker"
	"github.com/blokhinnv/gophermart/internal/app/models"

	"github.com/blokhinnv/gophermart/internal/app/server/config"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var STATUSES = map[string]int{
	"NEW":        0,
	"REGISTERED": 1,
	"PROCESSING": 2,
	"INVALID":    3,
	"PROCESSED":  4,
}

type DatabaseService struct {
	conn *pgxpool.Pool
}

func NewDatabaseService(
	cfg *config.Config,
	ctx context.Context,
	recreateOnStart bool,
) (*DatabaseService, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURI)
	if err != nil {
		return nil, err
	}
	conn, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}
	if recreateOnStart {
		log.Printf("Recreating tables in the DB...")
		if _, err = conn.Exec(ctx, createSQL); err != nil {
			return nil, err
		}
	}
	return &DatabaseService{conn: conn}, nil
}

func (db *DatabaseService) Tracker() ordertracker.Tracker {
	return ordertracker.NewDBTracker(db.conn)
}

func (db *DatabaseService) AddUser(
	ctx context.Context,
	username, pwd string,
) (*models.User, error) {
	log.Printf("Adding user %v...", username)
	salt, err := auth.GenerateSalt()
	if err != nil {
		return nil, err
	}
	pwdHash := auth.GenerateHash(pwd, salt)
	var addedID int
	err = db.conn.QueryRow(ctx, addUserSQL, username, pwdHash, salt).Scan(&addedID)
	if err != nil {
		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {
			if pgerr.Code == pgerrcode.UniqueViolation {
				return nil, fmt.Errorf("%w: %v", ErrUserAlreadyExists, username)
			}
		}
		return nil, err
	}
	return &models.User{ID: addedID, Username: username, HashedPassword: pwdHash, Salt: salt}, nil
}

func (db *DatabaseService) FindUser(
	ctx context.Context,
	username, pwd string,
) (*models.User, error) {
	var storedHash, salt string
	var id int
	err := db.conn.QueryRow(ctx, selectUserByLoginSQL, username).Scan(&id, &storedHash, &salt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %v", ErrUserNotFound, username)
		}
		return nil, err
	}
	return &models.User{ID: id, Username: username, HashedPassword: storedHash, Salt: salt}, nil

}

func (db *DatabaseService) FindOrderByID(
	ctx context.Context,
	orderID string,
) (*models.Order, error) {
	order := models.Order{}
	err := db.conn.QueryRow(ctx, selectOrderByIDSQL, orderID).
		Scan(&order.ID, &order.UserID, &order.StatusID, &order.UploadedAt)
	if err != nil {
		return nil, err
	}
	return &order, err
}

func (db *DatabaseService) AddOrder(
	ctx context.Context,
	orderID string,
	userID int,
) error {
	log.Printf("Adding order orderID=%v userID=%v...", orderID, userID)
	// пытаемся найти заказ в БД
	order, err := db.FindOrderByID(ctx, orderID)
	if err != nil {
		// не нашли заказ - надо добавить
		if errors.Is(err, pgx.ErrNoRows) {
			_, err = db.conn.Exec(ctx, addOrderSQL, orderID, userID)
			return err
		}
		// любая другая ошибка - плохо
		return err
	}
	// удалось что-то найти
	if order.UserID == userID {
		return fmt.Errorf(
			"%w: orderID=%v userID=%v",
			ErrOrderAlreadyAddedByThisUser,
			orderID,
			userID,
		)
	} else {
		return fmt.Errorf("%w: orderID=%v userID=%v", ErrOrderAlreadyAddedByOtherUser, orderID, userID)
	}
}

func (db *DatabaseService) UpdateOrderStatus(ctx context.Context, orderID, newStatus string) error {
	_, err := db.conn.Exec(ctx, updateOrderStatusSQL, newStatus, orderID)
	return err
}

func (db *DatabaseService) AddAccrualRecord(
	ctx context.Context,
	orderID string,
	sum float64,
) error {
	log.Printf("Adding accrual record orderID=%v sum=%v...", orderID, sum)
	_, err := db.conn.Exec(ctx, addTransactionSQL, orderID, sum, "ACCRUAL")
	return err
}

func (db *DatabaseService) FindOrdersByUserID(
	ctx context.Context,
	userID int,
) ([]models.Order, error) {
	orders := make([]models.Order, 0)
	rows, err := db.conn.Query(ctx, getOrdersByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		order := models.Order{}
		if err := rows.Scan(&order.ID, &order.Status, &order.Accrual, &order.UploadedAt); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, ErrEmptyResult
	}
	return orders, nil
}

func (db *DatabaseService) GetBalance(ctx context.Context, userID int) (*models.Balance, error) {
	balance := models.Balance{}
	err := db.conn.QueryRow(ctx, getBalanceSQL, userID).
		Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		return nil, err
	}
	return &balance, err
}

func (db *DatabaseService) AddWithdrawalRecord(
	ctx context.Context,
	orderID string,
	sum float64,
	userID int,
) error {
	log.Printf("Adding withdrawal record orderID=%v sum=%v...", orderID, sum)
	// если заказа не было - добавим
	err := db.AddOrder(ctx, orderID, userID)
	if err != nil {
		if !errors.Is(err, ErrOrderAlreadyAddedByOtherUser) &&
			!errors.Is(err, ErrOrderAlreadyAddedByThisUser) {
			return err
		}
	}
	_, err = db.conn.Exec(ctx, addTransactionSQL, orderID, sum, "WITHDRAWAL")
	if err != nil {
		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {
			if pgerr.Code == pgerrcode.ForeignKeyViolation {
				return fmt.Errorf("%w: %v", ErrMissingOrderID, orderID)
			}
		}
	}
	return err
}

func (db *DatabaseService) GetWithdrawals(
	ctx context.Context,
	userID int,
) ([]models.Withdrawal, error) {
	withdrawals := make([]models.Withdrawal, 0)
	rows, err := db.conn.Query(ctx, getWithdrawalsSQL, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		wd := models.Withdrawal{}
		if err := rows.Scan(&wd.Order, &wd.Sum, &wd.ProcessedAt); err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, wd)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(withdrawals) == 0 {
		return nil, ErrEmptyResult
	}
	return withdrawals, nil
}

func (db *DatabaseService) Close() {
	log.Println("Closing DB connection...")
	db.conn.Close()
}
