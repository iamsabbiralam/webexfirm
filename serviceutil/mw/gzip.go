package mw

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// Gzip Compression
type gzResp struct {
	io.Writer
	http.ResponseWriter
}

func (w gzResp) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func Gzip(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			handler.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		handler.ServeHTTP(gzResp{Writer: gz, ResponseWriter: w}, r)
	})
}
