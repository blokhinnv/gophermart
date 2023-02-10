package handlers

import (
	"github.com/blokhinnv/gophermart/internal/app/accrual"
	"github.com/blokhinnv/gophermart/internal/app/database"
	"github.com/blokhinnv/gophermart/internal/app/server/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

func NewRouter(db database.Service, cfg *config.Config) chi.Router {
	r := chi.NewRouter()
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
	postOrder := NewPostOrder(db, 10, 2, accrualService)
	getOrder := GetOrder{db: db}
	balance := Balance{db: db}
	withdraw := Withdraw{db: db}

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
			r.Post("/withdraw", withdraw.Handler)
		})
	})

	return r
}
