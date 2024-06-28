package server

import (
	"flag"
	"log"
	"net/http"

	"github.com/caarlos0/env/v11"
	"go.uber.org/zap"

	"github.com/DenisquaP/ya-metrics/internal/server/handlers"
)

type config struct {
	RunAddr string `env:"ADDRESS" envDefault:"localhost:8080"`
}

func NewConfig() (config, error) {
	var cfg config

	// Setting values by flags, if env not empty, using env
	flag.StringVar(&cfg.RunAddr, "a", "localhost:8080", "address and port to run server")

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
	router := handlers.InitRouter(suggared)

	if err := http.ListenAndServe(cfg.RunAddr, router); err != nil {
		suggared.Fatalw("Failed to start server", "error", err)
	}
}
