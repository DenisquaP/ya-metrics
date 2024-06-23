package handlers

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

func (h *Handler) createMetric(rw http.ResponseWriter, r *http.Request) {
	typeMetric := chi.URLParam(r, "type")
	nameMetric := chi.URLParam(r, "name")
	valueMetric := chi.URLParam(r, "value")

	if nameMetric == "" {
		http.Error(rw, "empty name", http.StatusNotFound)
		return
	}

	if err := h.Metrics.WriteMetric(nameMetric, typeMetric, valueMetric); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	
	rw.WriteHeader(http.StatusOK)
}

func (h *Handler) GetMetric(rw http.ResponseWriter, r *http.Request) {
	typeMet := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")

	val, err := h.Metrics.GetMetric(typeMet, name)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(val))
}

func (h *Handler) GetMetrics(rw http.ResponseWriter, r *http.Request) {
	metrics := h.Metrics.GetMetrics()

	metHTML := strings.Replace(HTMLMet, "{{metrics}}", metrics, -1)

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(metHTML))
}

var HTMLMet = `<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Metrics</title>
</head>

<body>
    Metrics:
    {{metrics}}
</body>

</html>`
