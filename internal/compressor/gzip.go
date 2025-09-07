package compressor

import (
	"compress/gzip"
	"io"
	"net/http"

	"github.com/domurdoc/shortener/internal/httputils"
)

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
	ok bool
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{w: w, zw: gzip.NewWriter(w)}
}

func (c *compressWriter) Write(p []byte) (int, error) {
	if c.ok {
		return c.zw.Write(p)
	}
	return c.w.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if httputils.HasContentType(
		c.w.Header(),
		httputils.ContentTypeJSON,
		httputils.ContentTypeTextPlain,
	) {
		c.ok = true
		httputils.SetContentEncoding(c.w.Header(), httputils.EncodingGZIP)
	}
	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Close() error {
	return c.zw.Close()
}

type compressReader struct {
	r  io.ReadCloser
	gr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return &compressReader{r: r, gr: gr}, nil
}

func (c *compressReader) Read(p []byte) (int, error) {
	return c.gr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.gr.Close()
}

func GZIPMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w
		if httputils.HasAcceptsEncoding(r.Header, httputils.EncodingGZIP) {
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		}
		if httputils.HasContentEncoding(r.Header, httputils.EncodingGZIP) {
			body, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer body.Close()
			r.Body = body
		}
		handler.ServeHTTP(ow, r)
	})
}
