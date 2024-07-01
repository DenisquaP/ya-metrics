package yametrics

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/DenisquaP/ya-metrics/pkg/models"
)

func (m *MemStorage) SaveMetricsToFile(wd string) error {
	file, err := os.OpenFile(filepath.Join(wd, m.FilePath), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() error {
		if err := file.Close(); err != nil {
			return err
		}

		return nil
	}()

	metrics, err := toJSON(m)
	if err != nil {
		return err
	}

	_, err = file.Write(metrics)
	if err != nil {
		return err
	}

	return nil
}

func (m *MemStorage) RestoreFromFile(wd string) error {
	metrics, err := os.ReadFile(filepath.Join(wd, m.FilePath))
	if err != nil {
		return err
	}

	mSlice := make([]models.Metrics, 19)

	err = json.Unmarshal(metrics, &mSlice)
	if err != nil {
		return err
	}

	for _, metric := range mSlice {
		switch metric.MType {
		case "gauge":
			m.Gauge[metric.ID] = *metric.Value
		case "counter":
			m.Counter[metric.ID] = *metric.Delta
		}
	}

	return nil
}

func toJSON(m *MemStorage) ([]byte, error) {
	metrics := make([]models.Metrics, 0, len(m.Gauge)+len(m.Counter))

	for k, v := range m.Gauge {
		var m models.Metrics

		m.ID = k
		m.MType = "gauge"
		m.Value = &v

		metrics = append(metrics, m)
	}

	for k, v := range m.Counter {
		var m models.Metrics

		m.ID = k
		m.MType = "counter"
		m.Delta = &v

		metrics = append(metrics, m)
	}

	return json.Marshal(metrics)
}
