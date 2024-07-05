package yametrics

import (
	"context"
	"fmt"
)

// MemStorage struct
type MemStorage struct {
	Gauge    map[string]float64
	Counter  map[string]int64
	FilePath string
}

func NewMemStorage(filepath string) *MemStorage {
	return &MemStorage{
		Gauge:    make(map[string]float64),
		Counter:  make(map[string]int64),
		FilePath: filepath,
	}
}

// Запись метрики типа Gauge
func (m *MemStorage) WriteGauge(ctx context.Context, name string, val float64) (float64, error) {
	m.Gauge[name] = val
	return m.Gauge[name], nil
}

// Запись метрики типа Counter
func (m *MemStorage) WriteCounter(ctx context.Context, name string, val int64) (int64, error) {
	m.Counter[name] += val
	return m.Counter[name], nil
}

// Получение всех метрик
func (m *MemStorage) GetMetrics(ctx context.Context) (string, error) {
	res := ""

	for k, v := range m.Gauge {
		res += fmt.Sprintf("%v: %v\n", k, v)
	}
	for k, v := range m.Counter {
		res += fmt.Sprintf("%v: %v\n", k, v)
	}

	return res, nil
}

// Получение метрики типа Gauge
func (m *MemStorage) GetGauge(ctx context.Context, name string) (float64, error) {
	g, ok := m.Gauge[name]
	if !ok {
		return 0, fmt.Errorf("variable does not exist")
	}

	return g, nil
}

// Получение метрики типа Counter
func (m *MemStorage) GetCounter(ctx context.Context, name string) (int64, error) {
	c, ok := m.Counter[name]
	if !ok {
		return 0, fmt.Errorf("variable does not exists")
	}

	return c, nil
}
