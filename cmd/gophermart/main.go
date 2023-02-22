package main

import (
	"log"
	"os"

	_ "github.com/blokhinnv/gophermart/internal/app"
	"github.com/blokhinnv/gophermart/internal/app/server"
	"github.com/blokhinnv/gophermart/internal/app/server/config"
	"github.com/joho/godotenv"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	godotenv.Load(".env")
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	server.RunServer(cfg)
}
