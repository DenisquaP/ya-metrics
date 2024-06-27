package middlewares

import (
	
	"net/http"
	"time"

	"go.uber.org/zap"
)

func Logging(logger zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//body, err := io.ReadAll(r.Body)
			//if err != nil {
				//logger.Error(err)
			//}
			//r.Body = io.NopCloser(bytes.NewBuffer(body))

			//var jsBody models.Metrics
			//if err := json.Unmarshal(body, &jsBody); err != nil {
				//logger.Error(err)
			//}

			ts := time.Now()

			lw := loggingResponseWriter{
				ResponseWriter: w,
				responseData: &responseData{
					status: 0,
					size:   0,
				}}

			next.ServeHTTP(&lw, r)

			// request logging
			logger.Infow("request", "method", r.Method, "url", r.URL, "time", time.Since(ts))

			//var jsBody2 models.Metrics
			//if err := json.Unmarshal(lw.responseData.body, &jsBody2); err != nil {
				//logger.Error(err)
			//}

			// response logging
			logger.Infow("response", "status", lw.responseData.status, "size", lw.responseData.size)
		})
	}
}

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.responseData.status = statusCode // захватываем код статуса
	r.ResponseWriter.WriteHeader(statusCode)
}
