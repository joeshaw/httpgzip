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
	gzipped       *gzip.Writer
	headerWritten bool
}

func (w *responseWriter) WriteHeader(status int) {
	w.headerWritten = true
	w.ResponseWriter.WriteHeader(status)
}

func (w *responseWriter) Write(data []byte) (int, error) {
	// Try to detect the content type if one wasn't provided.
	// If we don't do this, the underlying http.ResponseWriter
	// does and sets it to application/gzip.
	if !w.headerWritten && w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", http.DetectContentType(data))
	}

	return w.gzipped.Write(data)
}

func (w *responseWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		w.gzipped.Flush()
		f.Flush()
	}
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

		grw := responseWriter{
			ResponseWriter: w,
			gzipped:        gw,
		}

		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Content-Encoding", "gzip")

		h.ServeHTTP(&grw, r)
	})
}
