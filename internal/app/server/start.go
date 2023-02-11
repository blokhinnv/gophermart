package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/blokhinnv/gophermart/internal/app/database"
	"github.com/blokhinnv/gophermart/internal/app/server/config"
	"github.com/blokhinnv/gophermart/internal/app/server/handlers"
)

func RunServer(cfg *config.Config) {
	db, err := database.NewDatabaseService(cfg, context.Background(), true)
	if err != nil {
		log.Fatal(err)
	}

	shutdownCtx, _ := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	r := handlers.NewRouter(db, cfg)

	go func() {
		<-shutdownCtx.Done()
		log.Printf("Shutting down gracefully...")
		r.Shutdown()
		db.Close()
		log.Printf("Bye...")
		os.Exit(0)
	}()

	log.Printf("Starting server with config %+v\n", cfg)
	log.Fatal(http.ListenAndServe(cfg.RunAddress, r))
}
