package memyandex

import (
	"context"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v4/mem"
)

const MetricsUpdateURL = "http://%s/update/"
const AllMetricsURL = "http://%s/updates/"

type Counter int64
type Gauge float64

type MemStatsYaSt struct {
	RuntimeMem      *runtime.MemStats
	PollCount       int64
	RandomValue     float64
	TotalMemory     float64
	FreeMemory      float64
	CPUutilization1 float64
}

func (m *MemStatsYaSt) UpdateMetrics(ctx context.Context, errChan chan error) {
	runtime.ReadMemStats(m.RuntimeMem)
	m.RandomValue = float64(m.RuntimeMem.Alloc) / float64(1024)
	m.PollCount++

	v, err := mem.VirtualMemory()
	if err != nil {
		errChan <- err
	}

	m.TotalMemory = float64(v.Total)
	m.FreeMemory = float64(v.Free)

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
