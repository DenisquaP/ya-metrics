package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"

	"github.com/DenisquaP/ya-metrics/internal/server/middlewares"
	"github.com/DenisquaP/ya-metrics/internal/server/usecase"
)

type Handler struct {
	Metrics usecase.MetricService
	Logger  *zap.SugaredLogger
}

func NewHandler(metrics usecase.MetricService, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		Metrics: metrics,
		Logger:  logger,
	}
}

func NewRouterWithMiddlewares(ctx context.Context, logger *zap.SugaredLogger, metrics usecase.MetricService, key string) http.Handler {
	select {
	case <-ctx.Done():
		logger.Errorw("context canceled", "error", ctx.Err())
		return nil
	default:
	}

	r := chi.NewRouter()

	r.Use(middlewares.Logging(logger))

	// Middleware for compression
	r.Use(middlewares.Compression)

	h := NewHandler(metrics, logger)

	// To get all metrics in HTML
	r.Get("/", h.GetMetrics)

	// Ping database
	r.Get("/ping", h.Ping)

	r.Route("/", func(r chi.Router) {
		// Middleware for check ContentType
		r.Use(middleware.AllowContentType("application/json"))

		// Update metric
		r.Post("/update/{type}/{name}/{value}", h.createMetric)

		// Update metric JSON
		r.Post("/update/", h.createMetricJSON)

		// Get metric JSON
		r.Post("/value/", h.GetMetricJSON)
	})

	r.Route("/updates", func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json"))

		if key != "" {
			r.Use(middlewares.GetSum(logger, key))
		}

		// Update multiple metric
		r.Post("/", h.UpdateMultiple)
	})

	// Получение метрик v1
	r.Get("/value/{type}/{name}", h.GetMetric)

	return r
}
