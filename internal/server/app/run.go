package app

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/DenisquaP/ya-metrics/internal/server/config"
	"github.com/DenisquaP/ya-metrics/internal/server/db"
	"github.com/DenisquaP/ya-metrics/internal/server/handlers"
	yametrics "github.com/DenisquaP/ya-metrics/internal/server/yaMetrics"
	_ "github.com/DenisquaP/ya-metrics/migrations"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, err := newLogger()
	if err != nil {
		log.Fatal(err)
	}

	// Initiating config
	cfg, err := config.NewConfig()
	if err != nil {
		logger.Fatalw("Failed to parse config", "error", err)
	}

	logger.Infow("Starting server", "address", cfg.RunAddr)

	// Initiating router
	metrics := yametrics.NewMemStorage(cfg.FileStoragePath)

	// If database dsn is empty writing metrics to file and local struct
	var router http.Handler
	var database *db.DB
	if cfg.DatabaseDsn != "" {
		// Initiating DB
		database, err = db.NewDB(ctx, logger, cfg.DatabaseDsn)
		if err != nil {
			logger.Fatalw("Failed to create new DB", "error", err)
		}
		defer database.DB.Close(ctx)

		db, err := sql.Open("postgres", cfg.DatabaseDsn)
		if err != nil {
			logger.Fatalw("Failed to open DB", "error", err)
		}
		defer db.Close()

		if err := goose.Up(db, "./migrations"); err != nil {
			logger.Fatalw("Failed to run migrations", "error", err)
		}

		router = handlers.NewRouterWithMiddlewares(ctx, logger, database, cfg.Key)

	} else {
		router = handlers.NewRouterWithMiddlewares(ctx, logger, metrics, cfg.Key)

		go writeFile(ctx, logger, metrics, cfg)
	}

	if err := http.ListenAndServe(cfg.RunAddr, router); err != nil {
		logger.Fatalw("Failed to start server", "error", err)
	}
}
