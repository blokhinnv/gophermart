package handlers

import (
	"github.com/blokhinnv/gophermart/internal/app/database"
	"github.com/blokhinnv/gophermart/internal/app/server/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

func NewRouter(db *database.DatabaseService, cfg *config.Config) chi.Router {
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
	// TODO: DELETE
	_ = tokenAuth
	order := NewPostOrder(db, 10, 2, cfg.AccrualSystemAddress)

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
			r.Post("/orders", order.Handler)
		})
	})

	return r
}
