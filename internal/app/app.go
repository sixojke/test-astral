package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/sixojke/test-astral/internal/config"
	"github.com/sixojke/test-astral/internal/delivery"
	"github.com/sixojke/test-astral/internal/repository"
	"github.com/sixojke/test-astral/internal/server"
	"github.com/sixojke/test-astral/internal/service"
	"github.com/sixojke/test-astral/pkg/auth"
	"github.com/sixojke/test-astral/pkg/db"
	"github.com/sixojke/test-astral/pkg/hash"
	"github.com/sixojke/test-astral/pkg/logger"
	"github.com/sixojke/test-astral/pkg/migrations"
)

const (
	configs = "configs"
	env     = ".env"
)

// @title All social networks shop API
// @version 1.0
// @description REST API for shop

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey UsersAuth
// @in header
// @name Authorization
func Run() {
	// Get project directories
	currentDir, err := os.Getwd()
	if err != nil {
		logger.Fatalf("failed to get current directory: %v", err)
	}

	// Init config
	cfg, err := config.Init(path.Join(currentDir, configs), path.Join(currentDir, env))
	if err != nil {
		log.Fatal(err)
	}

	// Init logger
	enableLogger(cfg.Logger.LogLevel)

	// Init hasher
	hasher := hash.NewSHA1Hasher(cfg.Hasher.Salt)

	// Init token manager
	tokenManager, err := auth.NewManager(cfg.Authorization.JWT.SigningKey)
	if err != nil {
		logger.Fatalf("error init token manager: %v", err)
	}

	// Init PostgreSQL
	postgres, err := db.NewPostgresDB(db.PostgresConfig{
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		Username: cfg.Postgres.Username,
		Password: cfg.Postgres.Password,
		DBName:   cfg.Postgres.DBName,
		SSLMode:  cfg.Postgres.SSLMode,
	})
	if err != nil {
		logger.Fatalf("error connect postgres db: %v", err)
	}
	defer postgres.Close()
	logger.Info("[POSTGRES] Connection successful")

	if err := migrations.MigratePostgres(cfg.Postgres); err != nil {
		logger.Errorf("postgres migrate error: %v", err)
	}
	logger.Info("[POSTGRES] Migrate successful")

	repo := repository.NewService(&repository.Deps{
		Postgres: postgres,
	})

	service := service.NewService(&service.Deps{
		Repository:   repo,
		Config:       cfg,
		Hasher:       hasher,
		TokenManager: tokenManager,
	})

	handler := delivery.NewHandler(service, cfg, tokenManager)

	srv := server.NewServer(cfg.HTTPServer, handler.Init())
	go func() {
		if err := srv.Start(); !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occurred while running http server: %v\n", err)
		}
	}()
	logger.Infof("[SERVER] Started on port :%v", cfg.HTTPServer.Port)

	shutdown(srv, postgres)
}

func enableLogger(logLevel int) {
	logger.NewLogger(zerolog.Level(logLevel), os.Stdout)
}

func shutdown(srv *server.Server, postgres *sqlx.DB) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 3 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logger.Errorf("failed to stop server: %v", err)
	}

	postgres.Close()
}
