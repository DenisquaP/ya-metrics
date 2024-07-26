package middlewares

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Logging middleware for logging requests
func Logging(logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ts := time.Now()

			// response writer logging struct
			lw := loggingResponseWriter{
				ResponseWriter: w,
				responseData: &responseData{
					status: 0,
					size:   0,
				}}

			next.ServeHTTP(&lw, r)

			// request and response logging
			logger.Infow(
				"request and response",
				"method", r.Method,
				"url", r.URL,
				"time", time.Since(ts),
				"status", lw.responseData.status,
				"size", lw.responseData.size,
			)
		})
	}
}

type (
	// response data struct for logging
	responseData struct {
		status int
		size   int
	}

	// custom response writer
	loggingResponseWriter struct {
		http.ResponseWriter // embed http.ResponseWriter
		responseData        *responseData
	}
)

// Write to write response data to logs and ResponseWriter
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // getting response size
	return size, err
}

// WriteHeader to write response status code to custom ResponseWriter
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.responseData.status = statusCode // getting response status
	r.ResponseWriter.WriteHeader(statusCode)
}
