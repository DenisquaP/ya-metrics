package compression

import (
	"compress/gzip"
	"io"
)

type compressReader struct {
	r  io.ReadCloser
	gz *gzip.Reader
}

func NewCompressReader(r io.ReadCloser) (*compressReader, error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		gz: gz,
	}, nil
}

func (c *compressReader) Read(b []byte) (int, error) {
	return c.gz.Read(b)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}

	return c.gz.Close()
}
