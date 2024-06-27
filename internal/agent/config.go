package agent

import (
	"context"
	"flag"
	"log"
	"runtime"

	"github.com/DenisquaP/ya-metrics/internal/agent/memyandex"
	"github.com/caarlos0/env/v11"
)

type config struct {
	RunAddr string `env:"ADDRESS" envDefault:"localhost:8080"`

	// частота отправки метрик на сервер
	ReportInterval int `env:"REPORT_INTERVAL" envDefault:"10"`

	// частота опроса метрик из пакета runtime
	PollInterval int `env:"POLL_INTERVAL" envDefault:"2"`
}

func NewConfig() (config, error) {
	var cfg config

	// Setting values by flags, if env not empty, using env
	flag.StringVar(&cfg.RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&cfg.ReportInterval, "r", 10, "interval between report calls")
	flag.IntVar(&cfg.PollInterval, "p", 2, "interval between polling calls")

	if err := env.Parse(&cfg); err != nil {
		return config{}, err
	}

	flag.Parse()
	return cfg, nil
}

func Run() {
	cfg, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	mem := memyandex.MemStatsYaSt{RuntimeMem: &runtime.MemStats{}}

	ctx := context.Background()

	for {
		mem.UpdateMetrics(ctx, cfg.PollInterval)
		if err := mem.SendToServer(ctx, cfg.RunAddr, cfg.ReportInterval); err != nil {
			log.Printf("error send metrics: %s", err)
		}
	}
}
