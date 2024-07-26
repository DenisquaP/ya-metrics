package app

import (
	"context"
	"log"
	"runtime"
	"sync/atomic"
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
	errChan := make(chan error)

	var rateCount atomic.Int64

	// Loop update and send metrics
	for {
		select {
		case err = <-errChan:
			close(errChan)
			sugared.Fatalw("error in goroutine", "error", err)
		case <-ctx.Done():
			return
		case <-tickerSend.C:
			if int(rateCount.Load()) > cfg.RateLimit {
				// if count of requests to server > RateLimit from cfg, wait for a second
				<-time.After(1 * time.Second)
			}
			go mem.SendAllMetricsToServer(ctx, cfg.RunAddr, cfg.Key, errChan, &rateCount)
		case <-tickerUpdate.C:
			go mem.UpdateMetrics(ctx, errChan)
		}

	}
}
