package yametrics

import (
	"fmt"
)

type Metric interface {
	MetricGetter
	MetricWriter
}

type MetricWriter interface {
	WriteGouge(name string, val float64) error
	WriteCounter(name string, val int64) error
}

type MetricGetter interface {
	GetMetrics() string
	GetGauge(name string) (float64, error)
	GetCounter(name string) (int64, error)
}

// MemStorage struct
type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
}

// Запись метрики типа Gauge
func (m *MemStorage) WriteGauge(name string, val float64) (float64, error) {
	m.Gauge[name] = val
	return m.Gauge[name], nil
}

// Запись метрики типа Counter
func (m *MemStorage) WriteCounter(name string, val int64) (int64, error) {
	m.Counter[name] += val
	return m.Counter[name], nil
}

// Получение всех метрик
func (m *MemStorage) GetMetrics() string {
	res := ""

	for k, v := range m.Gauge {
		res += fmt.Sprintf("%v: %v\n", k, v)
	}
	for k, v := range m.Counter {
		res += fmt.Sprintf("%v: %v\n", k, v)
	}

	return res
}

// Получение метрики типа Gauge
func (m *MemStorage) GetGauge(name string) (float64, error) {
	g, ok := m.Gauge[name]
	if !ok {
		return 0, fmt.Errorf("variable does not exists")
	}

	return g, nil
}

// Получение метрики типа Counter
func (m *MemStorage) GetCounter(name string) (int64, error) {
	c, ok := m.Counter[name]
	if !ok {
		return 0, fmt.Errorf("variable does not exists")
	}

	return c, nil
}
