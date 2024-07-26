// Package for operate with middlewares

package middlewares

import (
	"bytes"
	"io"
	"net/http"

	"github.com/DenisquaP/ya-metrics/internal/cryptography"
	"go.uber.org/zap"
)

// GetSum middleware checks hash SHA256 in header and compares it with hash from body
func GetSum(logger *zap.SugaredLogger, key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// getting metrics from body
			metrics, err := io.ReadAll(r.Body)
			if err != nil {
				logger.Errorw("error reading body", "error", err)
			}

			// getting hash from header
			sumGet := r.Header.Get("HashSHA256")
			if sumGet == "" {
				logger.Warnw("Missing hash SHA256 header")
				w.WriteHeader(http.StatusBadRequest)

				return
			}

			// calculating hash SHA256 from body
			expectedSum := cryptography.GetSum(metrics, key)

			// comparing hashes
			// if not equal, return bad request
			if sumGet != expectedSum {
				logger.Errorw("Expected hash does not match", "expected", expectedSum, "actual", sumGet)

				w.WriteHeader(http.StatusBadRequest)
			}

			// writing metrics to body
			r.Body = io.NopCloser(bytes.NewBuffer(metrics))
			next.ServeHTTP(w, r)
		})
	}
}
