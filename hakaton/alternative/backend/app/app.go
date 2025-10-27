package app

import (
	"backend/clients"
	"backend/config"
	"backend/repository"
	"backend/router"
	"backend/store"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

type App struct {
	Repo         *repository.Repository
	RedisStore   *store.RedisStore
	MinIOStore   *store.MinIOStore
	PythonClient *clients.PythonServiceClient
	Router       *mux.Router
	Logger       zerolog.Logger
}

func NewApp(cfg *config.Config) (*App, error) {
	repo, err := repository.NewRepository(&cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize main repository: %w", err)
	}

	redisStore, err := store.NewRedisStore(&cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize redis store: %w", err)
	}

	minioStore, err := store.NewMinIOStore(&cfg.MinIO)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize minio store: %w", err)
	}

	pythonClient := clients.NewPythonServiceClient(cfg.PythonService.URL)

	r := router.NewRouter(redisStore, minioStore, repo, pythonClient, cfg)

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	return &App{
		Repo:         repo,
		RedisStore:   redisStore,
		MinIOStore:   minioStore,
		PythonClient: pythonClient,
		Router:       r,
		Logger:       logger,
	}, nil
}

func (a *App) Run(addr string) error {
	server := &http.Server{
		Addr:         addr,
		Handler:      a.Router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	a.Logger.Info().Str("addr", addr).Msg("Starting server")
	return server.ListenAndServe()
}
