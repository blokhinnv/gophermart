package server

import (
	"context"
	"log"
	"net/http"

	"github.com/blokhinnv/gophermart/internal/app/database"
	"github.com/blokhinnv/gophermart/internal/app/server/config"
	"github.com/blokhinnv/gophermart/internal/app/server/handlers"
)

func RunServer(cfg *config.Config) {
	db, err := database.NewDatabaseService(cfg, context.Background(), true)
	if err != nil {
		log.Fatal(err)
	}
	r := handlers.NewRouter(db, cfg)
	log.Printf("Starting server with config %+v\n", cfg)
	log.Fatal(http.ListenAndServe(cfg.RunAddress, r))
}
