package compression

import (
	"compress/gzip"
	"io"
)

type CompressReader struct {
	r  io.ReadCloser
	gz *gzip.Reader
}

func NewCompressReader(r io.ReadCloser) (*CompressReader, error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &CompressReader{
		r:  r,
		gz: gz,
	}, nil
}

func (c *CompressReader) Read(b []byte) (int, error) {
	return c.gz.Read(b)
}

func (c *CompressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}

	return c.gz.Close()
}
