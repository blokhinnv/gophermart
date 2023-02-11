package handlers

import (
	"context"

	"github.com/blokhinnv/gophermart/internal/app/accrual"
	"github.com/blokhinnv/gophermart/internal/app/database"
	"github.com/blokhinnv/gophermart/internal/app/server/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

type Router struct {
	*chi.Mux
	handlerContextCloseFuncs []context.CancelFunc
}

func (r *Router) Shutdown() {
	for _, cancel := range r.handlerContextCloseFuncs {
		cancel()
	}
}

func NewRouter(db database.Service, cfg *config.Config) Router {
	r := Router{
		Mux:                      chi.NewRouter(),
		handlerContextCloseFuncs: make([]context.CancelFunc, 0),
	}
	reg := Register{
		LogReg: LogReg{
			db:             db,
			signingKey:     []byte(cfg.JWTSigningKey),
			expireDuration: cfg.JWTExpireDuration,
		},
	}
	login := Login{
		LogReg: LogReg{
			db:             db,
			signingKey:     []byte(cfg.JWTSigningKey),
			expireDuration: cfg.JWTExpireDuration,
		},
	}
	tokenAuth := jwtauth.New("HS256", []byte(cfg.JWTSigningKey), nil)
	accrualService := accrual.NewAccrualService(cfg.AccrualSystemAddress)
	postOrder, cancel := NewPostOrder(db, 2, accrualService)
	r.handlerContextCloseFuncs = append(r.handlerContextCloseFuncs, cancel)
	getOrder := GetOrder{db: db}
	balance := Balance{db: db}
	withdraw := Withdraw{db: db}
	withdrawals := Withdrawals{db: db}

	r.Use(middleware.Logger)
	r.Route("/api/user", func(r chi.Router) {
		// доступны без авторизации
		r.Group(func(r chi.Router) {
			r.Post("/register", reg.Handler)
			r.Post("/login", login.Handler)
		})
		// доступны с авторизацией
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator)
			r.Post("/orders", postOrder.Handler)
			r.Get("/orders", getOrder.Handler)
			r.Get("/balance", balance.Handler)
			r.Post("/balance/withdraw", withdraw.Handler)
			r.Get("/withdrawals", withdrawals.Handler)
		})
	})

	return r
}
