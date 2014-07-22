// Package httpgzip provides an http.Handler wrapper that
// transparently compresses the response payload with gzip if the
// "Accept-Encoding: gzip" request header is provided.  It sets the
// "Vary: Accept-Encoding" and "Content-Encoding: gzip" response
// headers.
package httpgzip

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type responseWriter struct {
	http.ResponseWriter
	gzipped *gzip.Writer
}

func (w *responseWriter) Write(data []byte) (int, error) {
	return w.gzipped.Write(data)
}

// Handler wraps the provided http.Handler with one that provides
// transparent gzip content encoding.
func Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}

		gw := gzip.NewWriter(w)
		defer gw.Close()

		grw := responseWriter{w, gw}

		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Content-Encoding", "gzip")

		h.ServeHTTP(&grw, r)
	})
}
