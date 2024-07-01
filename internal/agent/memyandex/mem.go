package memyandex

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/DenisquaP/ya-metrics/pkg/models"
)

const MetricsUpdateURL = "http://%s/update/"

type Counter int64
type Gauge float64

type Sender interface {
	// Sending metrics to server
	Send(addr, name string) error
}

// Send sends counter metrics to server
func (c Counter) Send(addr, name string) error {
	intC := int64(c)
	req := models.Metrics{
		ID:    name,
		MType: "counter",
		Delta: &intC,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	// Creating a new gzip writer
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	defer gz.Close()

	// Writing the body to the gzip writer
	if _, err = gz.Write(body); err != nil {
		return err
	}

	if err = gz.Flush(); err != nil {
		return err
	}

	// Sending request with compressed data
	client := http.Client{Timeout: 20 * time.Second}
	reqw, err := http.NewRequest("POST", fmt.Sprintf(MetricsUpdateURL, addr), &buf)
	if err != nil {
		return err
	}
	reqw.Header.Set("Content-Type", "application/json")
	reqw.Header.Set("Content-Encoding", "gzip")
	reqw.Header.Set("Accept-Encoding", "gzip")

	resp, err := client.Do(reqw)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("not expected status code: %d", resp.StatusCode)
	}

	return nil
}

// Send sends gauge metrics to server
func (g Gauge) Send(addr, name string) error {
	floatG := float64(g)
	req := models.Metrics{
		ID:    name,
		MType: "gauge",
		Value: &floatG,
	}

	// Creating a new gzip writer
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	defer gz.Close()

	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	// Writing the body to the gzip writer
	if _, err := gz.Write(body); err != nil {
		return err
	}

	if err := gz.Flush(); err != nil {
		return err
	}

	client := http.Client{}

	// Sending request with compressed data
	reqw, err := http.NewRequest("POST", fmt.Sprintf(MetricsUpdateURL, addr), &buf)
	if err != nil {
		return err
	}
	reqw.Header.Set("Content-Type", "application/json")
	reqw.Header.Set("Content-Encoding", "gzip")
	reqw.Header.Set("Accept-Encoding", "gzip")

	resp, err := client.Do(reqw)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("not expected status code: %d", resp.StatusCode)
	}

	return nil
}

func Send(sender Sender, addr, name string) error {
	return sender.Send(addr, name)
}

type MemStatsYaSt struct {
	RuntimeMem  *runtime.MemStats
	PollCount   int64
	RandomValue float64
}

func (m *MemStatsYaSt) UpdateMetrics(ctx context.Context, pollInterval int) {
	runtime.ReadMemStats(m.RuntimeMem)
	m.RandomValue = float64(m.RuntimeMem.Alloc) / float64(1024)
	m.PollCount++

	withTimeout, cancel := context.WithTimeout(ctx, time.Duration(pollInterval)*time.Second)
	defer cancel()
	<-withTimeout.Done()
}

func (m *MemStatsYaSt) SendToServer(ctx context.Context, runAddr string, reportInterval int) error {
	// Map of gauge metrics
	fields := []struct {
		value Gauge
		name  string
	}{
		{Gauge(float64(m.RuntimeMem.Alloc)), "Alloc"},
		{Gauge(float64(m.RuntimeMem.BuckHashSys)), "BuckHashSys"},
		{Gauge(float64(m.RuntimeMem.Frees)), "Frees"},
		{Gauge(m.RuntimeMem.GCCPUFraction), "GCCPUFraction"},
		{Gauge(float64(m.RuntimeMem.GCSys)), "GCSys"},
		{Gauge(float64(m.RuntimeMem.HeapAlloc)), "HeapAlloc"},
		{Gauge(float64(m.RuntimeMem.HeapIdle)), "HeapIdle"},
		{Gauge(float64(m.RuntimeMem.HeapObjects)), "HeapObjects"},
		{Gauge(float64(m.RuntimeMem.HeapReleased)), "HeapReleased"},
		{Gauge(float64(m.RuntimeMem.HeapSys)), "HeapSys"},
		{Gauge(float64(m.RuntimeMem.LastGC)), "LastGC"},
		{Gauge(float64(m.RuntimeMem.Lookups)), "Lookups"},
		{Gauge(float64(m.RuntimeMem.MCacheInuse)), "MCacheInuse"},
		{Gauge(float64(m.RuntimeMem.MCacheSys)), "MCacheSys"},
		{Gauge(float64(m.RuntimeMem.MSpanInuse)), "MSpanInuse"},
		{Gauge(float64(m.RuntimeMem.MSpanSys)), "MSpanSys"},
		{Gauge(float64(m.RuntimeMem.Mallocs)), "Mallocs"},
		{Gauge(float64(m.RuntimeMem.NextGC)), "NextGC"},
		{Gauge(float64(m.RuntimeMem.NumForcedGC)), "NumForcedGC"},
		{Gauge(float64(m.RuntimeMem.NumGC)), "NumGC"},
		{Gauge(float64(m.RuntimeMem.OtherSys)), "OtherSys"},
		{Gauge(float64(m.RuntimeMem.PauseTotalNs)), "PauseTotalNs"},
		{Gauge(float64(m.RuntimeMem.StackInuse)), "StackInuse"},
		{Gauge(float64(m.RuntimeMem.Sys)), "Sys"},
		{Gauge(float64(m.RuntimeMem.StackSys)), "StackSys"},
		{Gauge(float64(m.RuntimeMem.TotalAlloc)), "TotalAlloc"},
		{Gauge(float64(m.RuntimeMem.HeapInuse)), "HeapInuse"},
		{Gauge(m.RandomValue), "RandomValue"},
	}

	// Sending gauge metrics
	for _, field := range fields {
		if err := Send(field.value, runAddr, field.name); err != nil {
			return err
		}
	}

	// Sending counter metric
	if err := Send(Counter(m.PollCount), runAddr, "PollCount"); err != nil {
		return err
	}

	WithTimeout, cancel := context.WithTimeout(ctx, time.Duration(reportInterval)*time.Second)
	defer cancel()

	<-WithTimeout.Done()
	return nil
}
