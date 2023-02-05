package database

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/blokhinnv/gophermart/internal/app/auth"
	"github.com/blokhinnv/gophermart/internal/app/models"

	"github.com/blokhinnv/gophermart/internal/app/server/config"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

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

func (db *DatabaseService) GetConn() *pgxpool.Pool {
	return db.conn
}

func (db *DatabaseService) AddUser(ctx context.Context, username, pwd string) error {
	log.Printf("Adding user %v...", username)
	salt, err := auth.GenerateSalt()
	if err != nil {
		return nil
	}
	pwdHash := auth.GenerateHash(pwd, salt)
	_, err = db.conn.Exec(ctx, addUserSQL, username, pwdHash, salt)
	if err != nil {
		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {
			if pgerr.Code == pgerrcode.UniqueViolation {
				return fmt.Errorf("%w: %v", ErrUserAlreadyExists, username)
			}
		}
		return err
	}
	return nil
}

func (db *DatabaseService) FindUser(
	ctx context.Context,
	username, pwd string,
) (*models.User, error) {
	var storedHash, salt string
	err := db.conn.QueryRow(ctx, selectUserByLoginSQL, username).Scan(&storedHash, &salt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %v", ErrUserNotFound, username)
		}
		return nil, err
	}
	return &models.User{Username: username, HashedPassword: storedHash, Salt: salt}, nil

}
