package handlers

import (
	"fmt"
	"net/http"
)

// Ping pings db
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := h.Metrics.Ping(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		h.Logger.Errorw("ping error", "error", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status":"ok"}`)
}
