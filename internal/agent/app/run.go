package app

import (
	"context"
	"log"
	"runtime"
	"time"

	"github.com/DenisquaP/ya-metrics/internal/agent/config"
	"github.com/DenisquaP/ya-metrics/internal/agent/memyandex"
	"go.uber.org/zap"
)

func Run() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	sugared := logger.Sugar()

	cfg, err := config.NewConfig()
	if err != nil {
		sugared.Fatalw("Failed to parse config", "error", err)
	}

	// Creating struct for collecting metrics
	mem := memyandex.MemStatsYaSt{RuntimeMem: &runtime.MemStats{}}

	ctx := context.Background()

	// Tickers to send and update metrics
	tickerSend := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)
	tickerUpdate := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)

	// Loop update and send metrics
	for {
		select {
		case <-ctx.Done():
			return
		case <-tickerSend.C:
			if err := mem.SendAllMetricsToServer(ctx, cfg.RunAddr, cfg.Key); err != nil {
				sugared.Errorw("Failed to send metrics", "error", err)
			}

		case <-tickerUpdate.C:
			mem.UpdateMetrics(ctx)
		}

	}
}
