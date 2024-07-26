package usecase

import (
	"go.uber.org/zap"
)

type Metric struct {
	m      MetricService // struct implementation interfacep
	logger *zap.SugaredLogger
}

func NewMetric(m MetricService, logger *zap.SugaredLogger) *Metric {
	return &Metric{
		m:      m,
		logger: logger,
	}
}
