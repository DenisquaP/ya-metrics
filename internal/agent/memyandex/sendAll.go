package memyandex

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DenisquaP/ya-metrics/internal/cryptography"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/DenisquaP/ya-metrics/internal/agent/compress"
	"github.com/DenisquaP/ya-metrics/internal/models"
	"github.com/DenisquaP/ya-metrics/internal/repeat"
)

// Funcs to get pointer
func getPointerFloat(v float64) *float64 {
	return &v
}

func getPointerInt(v int64) *int64 {
	return &v
}

// SendAllMetricsToServer sends all metrics to server
func (m *MemStatsYaSt) SendAllMetricsToServer(ctx context.Context, addr string, key string, errChan chan error, rateCount *atomic.Int64) {
	// counting actual requests
	rateCount.Add(1)
	defer rateCount.Add(-1)

	// Metrics slice
	met := m.getSliceMetrics()

	metrics, err := json.Marshal(met)
	if err != nil {
		errChan <- err
		return
	}

	// Getting compressed data
	buf, err := compress.GetGZip(metrics)
	if err != nil {
		errChan <- err
		return
	}

	// Sending request with compressed data
	client := http.Client{Timeout: 5 * time.Second}
	reqw, err := http.NewRequest("POST", fmt.Sprintf(AllMetricsURL, addr), buf)
	if err != nil {
		errChan <- err
		return
	}
	reqw.Header.Set("Content-Type", "application/json")
	reqw.Header.Set("Content-Encoding", "gzip")
	reqw.Header.Set("Accept-Encoding", "gzip")

	// if key not nil writing sum to header
	if key != "" {
		sum := cryptography.GetSum(metrics, key)
		reqw.Header.Set("HashSHA256", sum)
	}

	resp, err := client.Do(reqw)
	if err != nil {
		var urlErr *net.OpError

		// Check if error is OpError
		if errors.As(err, &urlErr) {
			if err := repeat.RepeatNet(ctx, &client, reqw); err != nil {
				errChan <- err
				return
			}
		} else {
			errChan <- err
			return
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errChan <- fmt.Errorf("not expected status code: %d", resp.StatusCode)
		return
	}
}

// getSliceMetrics returns slice of metrics to send to server
func (m *MemStatsYaSt) getSliceMetrics() []models.Metrics {
	metrics := []models.Metrics{
		{
			ID:    "Alloc",
			MType: "gauge",
			Value: getPointerFloat(float64(m.RuntimeMem.Alloc)),
		},
		{
			ID:    "BuckHashSys",
			MType: "gauge",
			Value: getPointerFloat(float64(m.RuntimeMem.BuckHashSys)),
		},
		{
			ID:    "Frees",
			MType: "gauge",
			Value: getPointerFloat(float64(m.RuntimeMem.Frees)),
		},
		{
			ID:    "GCCPUFraction",
			MType: "gauge",
			Value: getPointerFloat(m.RuntimeMem.GCCPUFraction),
		},
		{
			ID:    "GCSys",
			MType: "gauge",
			Value: getPointerFloat(float64(m.RuntimeMem.GCSys)),
		},
		{
			ID:    "HeapAlloc",
			MType: "gauge",
			Value: getPointerFloat(float64(m.RuntimeMem.HeapAlloc)),
		},
		{
			ID:    "HeapIdle",
			MType: "gauge",
			Value: getPointerFloat(float64(m.RuntimeMem.HeapIdle)),
		},
		{
			ID:    "HeapObjects",
			MType: "gauge",
			Value: getPointerFloat(float64(m.RuntimeMem.HeapObjects)),
		},
		{
			ID:    "HeapReleased",
			MType: "gauge",
			Value: getPointerFloat(float64(m.RuntimeMem.HeapReleased)),
		},
		{
			ID:    "HeapSys",
			MType: "gauge",
			Value: getPointerFloat(float64(m.RuntimeMem.HeapSys)),
		},
		{
			ID:    "LastGC",
			MType: "gauge",
			Value: getPointerFloat(float64(m.RuntimeMem.LastGC)),
		},
		{
			ID:    "Lookups",
			MType: "gauge",
			Value: getPointerFloat(float64(m.RuntimeMem.Lookups)),
		},
		{
			ID:    "MCacheInuse",
			MType: "gauge",
			Value: getPointerFloat(float64(m.RuntimeMem.MCacheInuse)),
		},
		{
			ID:    "MCacheSys",
			MType: "gauge",
			Value: getPointerFloat(float64(m.RuntimeMem.MCacheSys)),
		},
		{
			ID:    "MSpanInuse",
			MType: "gauge",
			Value: getPointerFloat(float64(m.RuntimeMem.MSpanInuse)),
		},
		{
			ID:    "MSpanSys",
			MType: "gauge",
			Value: getPointerFloat(float64(m.RuntimeMem.MSpanSys)),
		},
		{
			ID:    "Mallocs",
			MType: "gauge",
			Value: getPointerFloat(float64(m.RuntimeMem.Mallocs)),
		},
		{
			ID:    "NextGC",
			MType: "gauge",
			Value: getPointerFloat(float64(m.RuntimeMem.NextGC)),
		},
		{
			ID:    "NumForcedGC",
			MType: "gauge",
			Value: getPointerFloat(float64(m.RuntimeMem.NumForcedGC)),
		},
		{
			ID:    "NumGC",
			MType: "gauge",
			Value: getPointerFloat(float64(m.RuntimeMem.NumGC)),
		},
		{
			ID:    "OtherSys",
			MType: "gauge",
			Value: getPointerFloat(float64(m.RuntimeMem.OtherSys)),
		},
		{
			ID:    "PauseTotalNs",
			MType: "gauge",
			Value: getPointerFloat(float64(m.RuntimeMem.PauseTotalNs)),
		},
		{
			ID:    "StackInuse",
			MType: "gauge",
			Value: getPointerFloat(float64(m.RuntimeMem.StackInuse)),
		},
		{
			ID:    "Sys",
			MType: "gauge",
			Value: getPointerFloat(float64(m.RuntimeMem.Sys)),
		},
		{
			ID:    "StackSys",
			MType: "gauge",
			Value: getPointerFloat(float64(m.RuntimeMem.StackSys)),
		},
		{
			ID:    "TotalAlloc",
			MType: "gauge",
			Value: getPointerFloat(float64(m.RuntimeMem.TotalAlloc)),
		},
		{
			ID:    "HeapInuse",
			MType: "gauge",
			Value: getPointerFloat(float64(m.RuntimeMem.HeapInuse)),
		},
		{
			ID:    "RandomValue",
			MType: "gauge",
			Value: getPointerFloat(m.RandomValue),
		},
		{
			ID:    "TotalMemory",
			MType: "gauge",
			Value: getPointerFloat(m.TotalMemory),
		},
		{
			ID:    "FreeMemory",
			MType: "gauge",
			Value: getPointerFloat(m.FreeMemory),
		},
		{
			ID:    "PollCount",
			MType: "counter",
			Delta: getPointerInt(m.PollCount),
		},
	}

	return metrics
}
