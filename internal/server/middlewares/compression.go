package middlewares

import (
	"net/http"
	"strings"

	"github.com/DenisquaP/ya-metrics/internal/server/compression"
)

// Compression middleware for compressing data
func Compression(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := w

		encodings := r.Header.Get("Accept-Encoding")
		// if content encoding contains gzip
		// using compress writer
		if strings.Contains(encodings, "gzip") {
			cw := compression.NewCompressWriter(w)
			rw = cw

			defer cw.Close()
		}

		// if request contains gzip encoding
		// using compress reader
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			cr, err := compression.NewCompressReader(r.Body)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}

			// setting compress reader to body
			r.Body = cr
			defer cr.Close()
		}

		next.ServeHTTP(rw, r)
	})
}
