package handlers

import (
	"github.com/blokhinnv/gophermart/internal/app/database"
	"github.com/blokhinnv/gophermart/internal/app/server/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	r.Use(middleware.Logger)
	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", reg.Handler)
		r.Post("/login", login.Handler)
	})
	return r
}
