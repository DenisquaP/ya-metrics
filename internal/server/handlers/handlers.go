package handlers

import (
	"net/http"

	"github.com/DenisquaP/ya-metrics/internal/server/middlewares"
	yametrics "github.com/DenisquaP/ya-metrics/internal/server/yaMetrics"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

type Handler struct {
	Metrics *yametrics.MemStorage
}

func NewHandler() *Handler {
	metrics := yametrics.NewMemStorage()
	return &Handler{
		Metrics: metrics,
	}
}

func InitRouter(logger zap.SugaredLogger) http.Handler {
	r := chi.NewRouter()
	r.Use(middlewares.Logging(logger))

	h := NewHandler()

	// Получение всех метрик в HTML
	r.Get("/", h.GetMetrics)

	r.Route("/", func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json"))

		// Обновление метрик v1
		r.Post("/update/{type}/{name}/{value}", h.createMetric)

		// Обновление метрик v2
		r.Post("/update/", h.createMetricV2)

		// Получение метрик v2
		r.Post("/value/", h.GetMetricV2)
	})

	// Получение метрик v1
	r.Get("/value/{type}/{name}", h.GetMetric)

	return r
}
