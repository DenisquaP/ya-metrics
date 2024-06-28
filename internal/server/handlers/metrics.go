package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"

	"github.com/DenisquaP/ya-metrics/pkg/models"
)

func (h *Handler) createMetric(rw http.ResponseWriter, r *http.Request) {
	typeMetric := chi.URLParam(r, "type")
	nameMetric := chi.URLParam(r, "name")
	valueMetric := chi.URLParam(r, "value")

	if nameMetric == "" {
		http.Error(rw, "empty name", http.StatusNotFound)
		return
	}

	switch typeMetric {
	case "counter":
		val, err := strconv.ParseInt(valueMetric, 10, 64)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		if _, err := h.Metrics.WriteCounter(nameMetric, val); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

	case "gauge":
		val, err := strconv.ParseFloat(valueMetric, 64)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		if _, err := h.Metrics.WriteGauge(nameMetric, val); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
	default:
		http.Error(rw, "wrong type", http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (h *Handler) GetMetric(rw http.ResponseWriter, r *http.Request) {
	typeMet := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")

	var resp []byte

	switch typeMet {
	case "counter":
		val, err := h.Metrics.GetCounter(name)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		}

		resp = []byte(strconv.FormatInt(val, 10))
	case "gauge":
		val, err := h.Metrics.GetGauge(name)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		}

		resp = []byte(strconv.FormatFloat(val, 'f', -1, 64))
	default:
		http.Error(rw, "wrong type", http.StatusBadRequest)
		return
	}

	log.Println(string(resp))

	rw.WriteHeader(http.StatusOK)
	if _, err := rw.Write(resp); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) createMetricV2(rw http.ResponseWriter, r *http.Request) {
	var request models.Metrics

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if request.ID == "" {
		http.Error(rw, "empty name", http.StatusNotFound)
		return
	}

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

	resp, err := json.Marshal(request)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	if _, err = rw.Write(resp); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetMetricV2(rw http.ResponseWriter, r *http.Request) {
	var request models.Metrics
	var response models.Metrics

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	switch request.MType {
	case "counter":
		c, err := h.Metrics.GetCounter(request.ID)
		if err != nil {
			log.Println(err.Error() + "[c]")
			http.Error(rw, err.Error()+"not found counter", http.StatusNotFound)
			return
		}
		response.ID = request.ID
		response.MType = request.MType
		response.Delta = &c
	case "gauge":
		g, err := h.Metrics.GetGauge(request.ID)
		if err != nil {
			log.Println(err.Error() + "[g]")
			http.Error(rw, err.Error()+"not found gauge", http.StatusNotFound)
			return
		}
		response.ID = request.ID
		response.MType = request.MType
		response.Value = &g
	default:
		http.Error(rw, "wrong type", http.StatusBadRequest)
		return
	}

	res, err := json.Marshal(response)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	if _, err = rw.Write(res); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

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
