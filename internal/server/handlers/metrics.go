package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"

	"github.com/DenisquaP/ya-metrics/pkg/models"
)

// Create metric with query params
func (h *Handler) createMetric(rw http.ResponseWriter, r *http.Request) {
	// Getting params
	typeMetric := chi.URLParam(r, "type")
	nameMetric := chi.URLParam(r, "name")
	valueMetric := chi.URLParam(r, "value")

	if nameMetric == "" {
		http.Error(rw, "empty name", http.StatusNotFound)
		return
	}

	// Getting function from map for writing metric
	funcWrite, ok := metricWrite[typeMetric]
	if !ok {
		http.Error(rw, "wrong type", http.StatusBadRequest)
		return
	}

	// Writing metric into structure
	err := funcWrite(h.Metrics, nameMetric, valueMetric)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

// Get metric with query params
func (h *Handler) GetMetric(rw http.ResponseWriter, r *http.Request) {
	typeMet := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")

	var resp []byte

	// Getting function from map for getting metric
	funcGet, ok := metricGet[typeMet]
	if !ok {
		http.Error(rw, "wrong type", http.StatusBadRequest)
		return
	}

	// Getting metric
	metric, err := funcGet(h.Metrics, name)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	resp = []byte(metric)

	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusOK)
	if _, err := rw.Write(resp); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Create metric with json body
func (h *Handler) createMetricJSON(rw http.ResponseWriter, r *http.Request) {
	var request models.Metrics

	// Decoding json
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if request.ID == "" {
		http.Error(rw, "empty name", http.StatusNotFound)
		return
	}

	// Writing metric
	switch request.MType {
	case "counter":
		newVal, err := h.Metrics.WriteCounter(request.ID, *request.Delta)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		request.Delta = &newVal

	case "gauge":
		newVal, err := h.Metrics.WriteGauge(request.ID, *request.Value)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		request.Value = &newVal
	default:
		http.Error(rw, "wrong type", http.StatusBadRequest)
		return
	}

	// Encoding json
	resp, err := json.Marshal(request)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// Writing response
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	if _, err = rw.Write(resp); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Get metric with json body
func (h *Handler) GetMetricJSON(rw http.ResponseWriter, r *http.Request) {
	var metric models.Metrics

	// Decoding json
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Writing metric
	switch metric.MType {
	case "counter":
		c, err := h.Metrics.GetCounter(metric.ID)
		if err != nil {
			log.Println(err.Error() + "[c]")
			http.Error(rw, err.Error()+"not found counter", http.StatusNotFound)
			return
		}

		// Writing response struct
		metric.Delta = &c
	case "gauge":
		g, err := h.Metrics.GetGauge(metric.ID)
		if err != nil {
			log.Println(err.Error() + "[g]")
			http.Error(rw, err.Error()+"not found gauge", http.StatusNotFound)
			return
		}

		// Writing response struct
		metric.Value = &g
	default:
		http.Error(rw, "wrong type", http.StatusBadRequest)
		return
	}

	// Encoding json
	res, err := json.Marshal(metric)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// Writing response
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	if _, err = rw.Write(res); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Get all metrics
func (h *Handler) GetMetrics(rw http.ResponseWriter, r *http.Request) {
	metrics := h.Metrics.GetMetrics()

	metHTML := strings.Replace(HTMLMet, "{{metrics}}", metrics, -1)

	rw.Header().Set("Content-Type", "text/html")
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
