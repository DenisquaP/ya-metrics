package server

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/caarlos0/env/v11"
	"go.uber.org/zap"

	"github.com/DenisquaP/ya-metrics/internal/server/handlers"
	yametrics "github.com/DenisquaP/ya-metrics/internal/server/yaMetrics"
)

type config struct {
	// Server address and port
	RunAddr string `env:"ADDRESS" envDefault:"localhost:8080"`
	// Interval between saving metrics to file
	StoreInterval int `env:"STORE_INTERVAL" envDefault:"300"`
	// Path to file with metrics
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"/tmp/metrics-db.json"`
	// Restore metrics from file
	Restore bool `env:"RESTORE" envDefault:"true"`
}

func NewConfig() (config, error) {
	var cfg config

	// Setting values by flags, if env not empty, using env
	flag.StringVar(&cfg.RunAddr, "a", "localhost:8080", "address and port to run server")
	// Setting values by flags, if env not empty, using env
	flag.IntVar(&cfg.StoreInterval, "i", 300, "interval between saving metrics to file")
	// Setting values by flags, if env not empty, using env
	flag.StringVar(&cfg.FileStoragePath, "f", "/tmp/metrics-db.json", "path to file with metrics")
	// Setting values by flags, if env not empty, using env
	flag.BoolVar(&cfg.Restore, "r", true, "restore metrics from file")

	if err := env.Parse(&cfg); err != nil {
		return config{}, err
	}

	flag.Parse()
	return cfg, nil
}

func Run() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	suggared := *logger.Sugar()

	cfg, err := NewConfig()
	if err != nil {
		suggared.Fatalw("Failed to parse config", "error", err)
	}

	suggared.Infow("Starting server", "address", cfg.RunAddr)

	metrics := yametrics.NewMemStorage(cfg.FileStoragePath)

	router := handlers.InitRouter(suggared, metrics)

	go func() {
		withTimeout, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.StoreInterval)*time.Second)
		defer cancel()
		for {
			select {
			case <-withTimeout.Done():
				logger.Info("Save metrics to file")
				metrics.SaveToFile()
				return
			default:
				for {
					if err := metrics.SaveToFile(); err != nil {
						suggared.Errorw("Failed to save metrics to file", "error", err)
					}
				}
			}
		}
	}()

	if err := http.ListenAndServe(cfg.RunAddr, router); err != nil {
		suggared.Fatalw("Failed to start server", "error", err)
	}
}
