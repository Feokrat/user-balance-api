package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Feokrat/user-balance-api/internal/server"

	"github.com/Feokrat/user-balance-api/internal/delivery/http"
	"github.com/Feokrat/user-balance-api/internal/service"

	"github.com/Feokrat/user-balance-api/internal/repository"

	"github.com/Feokrat/user-balance-api/internal/database"

	"github.com/Feokrat/user-balance-api/internal/config"
)

const configFile = "configs/config"

func main() {
	logger := log.New(os.Stdout, "logger: ", log.Lshortfile)
	cfg, err := config.Init(configFile, logger)
	if err != nil {
		logger.Fatalf("failed to load application configuration: %s", err)
	}

	db, err := database.NewPostgresDB(cfg.Postgresql, logger)
	if err != nil {
		logger.Fatalf("error with database: %s", err)
	}
	defer database.ClosePostgresDB(db)

	repos := repository.NewRepositories(db, logger)
	services := service.NewServices(repos, logger)

	handlers := http.NewHandler(services, logger)

	server := server.NewHTTPserver(cfg, handlers.Init())
	go func() {
		if err := server.Run(); err != nil {
			logger.Printf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	logger.Printf("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Println("Gracefully shutting down...")

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := server.Stop(ctx); err != nil {
		logger.Printf("error occurred on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logger.Printf("error occurred on db connection close: %s", err.Error())
	}
}
