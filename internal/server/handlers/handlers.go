package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"

	"github.com/DenisquaP/ya-metrics/internal/server/middlewares"
	yametrics "github.com/DenisquaP/ya-metrics/internal/server/yaMetrics"
)

type Handler struct {
	Metrics *yametrics.MemStorage
}

func NewHandler(metrics *yametrics.MemStorage) *Handler {
	return &Handler{
		Metrics: metrics,
	}
}

func NewRouterWithMiddlewares(logger *zap.SugaredLogger, metrics *yametrics.MemStorage) http.Handler {
	r := chi.NewRouter()
	r.Use(middlewares.Logging(logger))

	// Middleware для сжатия
	r.Use(middlewares.Compression)

	h := NewHandler(metrics)

	// Получение всех метрик в HTML
	r.Get("/", h.GetMetrics)

	r.Route("/", func(r chi.Router) {
		// Middleware для проверки ContentType
		r.Use(middleware.AllowContentType("application/json"))

		// Обновление метрик v1
		r.Post("/update/{type}/{name}/{value}", h.createMetric)

		// Обновление метрик v2
		r.Post("/update/", h.createMetricJSON)

		// Получение метрик v2
		r.Post("/value/", h.GetMetricJSON)
	})

	// Получение метрик v1
	r.Get("/value/{type}/{name}", h.GetMetric)

	return r
}
