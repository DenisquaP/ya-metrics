package compression

import (
	"compress/gzip"
	"net/http"
)

type CompressWriter struct {
	w  http.ResponseWriter
	gz *gzip.Writer
}

func NewCompressWriter(w http.ResponseWriter) *CompressWriter {
	return &CompressWriter{
		w:  w,
		gz: gzip.NewWriter(w),
	}
}

func (c *CompressWriter) Write(p []byte) (int, error) {
	return c.gz.Write(p)
}

func (c *CompressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *CompressWriter) WriteHeader(statusCode int) {
	if statusCode < http.StatusMultipleChoices {
		c.w.Header().Set("Content-Encoding", "gzip")
	}

	c.w.WriteHeader(statusCode)
}

func (c CompressWriter) Close() error {
	return c.gz.Close()
}
