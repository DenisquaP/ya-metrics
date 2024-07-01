package handlers

import (
	"fmt"
	"strconv"

	yametrics "github.com/DenisquaP/ya-metrics/internal/server/yaMetrics"
)

// Mapping metric write
var metricWrite map[string]func(metric *yametrics.MemStorage, name, value string) error = map[string]func(metric *yametrics.MemStorage, name, value string) error{
	"counter": func(metric *yametrics.MemStorage, name, value string) error {
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}

		_, err = metric.WriteCounter(name, val)
		if err != nil {
			return err
		}

		return nil
	},
	"gauge": func(metric *yametrics.MemStorage, name, value string) error {
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}

		_, err = metric.WriteGauge(name, val)
		if err != nil {
			return err
		}

		return nil
	},
}

// Mapping metric get
var metricGet map[string]func(metric *yametrics.MemStorage, name string) (string, error) = map[string]func(metric *yametrics.MemStorage, name string) (string, error){
	"counter": func(metric *yametrics.MemStorage, name string) (string, error) {
		val, err := metric.GetCounter(name)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%v", val), nil
	},
	"gauge": func(metric *yametrics.MemStorage, name string) (string, error) {
		val, err := metric.GetGauge(name)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%v", val), nil
	},
}
