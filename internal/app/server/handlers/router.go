package handlers

import (
	"context"
	"log"

	"github.com/blokhinnv/gophermart/internal/app/accrual"
	"github.com/blokhinnv/gophermart/internal/app/database"
	"github.com/blokhinnv/gophermart/internal/app/server/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

type Router struct {
	*chi.Mux
	reg         *Register
	login       *Login
	postOrder   *PostOrder
	getOrder    *GetOrder
	balance     *Balance
	withdraw    *Withdraw
	withdrawals *Withdrawals
}

func (r *Router) Shutdown() {
	log.Println("Shutting down Router...")
	// хочу убедиться, что все горутины этого хендлера завершились,
	// прежде чем двигаться дальше
	r.postOrder.WaitDone()
}

func NewRouter(db database.Service, cfg *config.Config, serverCtx context.Context) Router {
	rt := Router{
		Mux: chi.NewRouter(),
	}
	rt.reg = &Register{
		LogReg: LogReg{
			db:             db,
			signingKey:     []byte(cfg.JWTSigningKey),
			expireDuration: cfg.JWTExpireDuration,
		},
	}
	rt.login = &Login{
		LogReg: LogReg{
			db:             db,
			signingKey:     []byte(cfg.JWTSigningKey),
			expireDuration: cfg.JWTExpireDuration,
		},
	}
	tokenAuth := jwtauth.New("HS256", []byte(cfg.JWTSigningKey), nil)
	accrualService := accrual.NewAccrualService(cfg.AccrualSystemAddress)
	rt.postOrder = NewPostOrder(
		db,
		2,
		serverCtx,
		accrualService,
		cfg.AccrualSystemPoolInterval,
	)

	rt.getOrder = &GetOrder{db: db}
	rt.balance = &Balance{db: db}
	rt.withdraw = &Withdraw{db: db}
	rt.withdrawals = &Withdrawals{db: db}

	rt.Use(middleware.Logger)
	rt.Route("/api/user", func(r chi.Router) {
		// доступны без авторизации
		r.Group(func(r chi.Router) {
			r.Post("/register", rt.reg.Handler)
			r.Post("/login", rt.login.Handler)
		})
		// доступны с авторизацией
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator)
			r.Post("/orders", rt.postOrder.Handler)
			r.Get("/orders", rt.getOrder.Handler)
			r.Get("/balance", rt.balance.Handler)
			r.Post("/balance/withdraw", rt.withdraw.Handler)
			r.Get("/withdrawals", rt.withdrawals.Handler)
		})
	})

	return rt
}
