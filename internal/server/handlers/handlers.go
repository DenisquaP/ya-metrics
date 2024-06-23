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

	r.Get("/", h.GetMetrics)

	r.Route("/update", func(r chi.Router) {
		r.Use(middleware.AllowContentType("text/plain"))

		r.Post("/{type}/{name}/{value}", h.createMetric)
	})

	r.Get("/value/{type}/{name}", h.GetMetric)

	return r
}
