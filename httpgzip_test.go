package httpgzip

import (
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

const lorem = `Lorem ipsum dolor sit amet, consectetur adipisicing
elit, sed do eiusmod tempor incididunt ut labore et dolore magna
aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco
laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor
in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla
pariatur. Excepteur sint occaecat cupidatat non proident, sunt in
culpa qui officia deserunt mollit anim id est laborum.`

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(lorem))
}

func TestHTTPGzip(t *testing.T) {
	handler := http.HandlerFunc(handlerFunc)
	wrapped := Handler(handler)

	// A call against the regular handler
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, &http.Request{})

	// A call against the wrapped handler without gzip
	rec2 := httptest.NewRecorder()
	wrapped.ServeHTTP(rec2, &http.Request{})

	if rec.Code != rec2.Code {
		t.Fatalf("Expected status code (%d) != actual (%d)", rec.Code, rec2.Code)
	}

	if !reflect.DeepEqual(rec.HeaderMap, rec2.HeaderMap) {
		t.Fatalf("Expected header map (%#v) != actual (%#v)", rec.HeaderMap, rec2.HeaderMap)
	}

	if rec.Body.String() != rec2.Body.String() {
		t.Fatal("Expected body != actual")
	}

	// A call against the wrapped handler with gzip on
	rec3 := httptest.NewRecorder()
	req := &http.Request{
		Header: map[string][]string{},
	}
	req.Header.Set("Accept-Encoding", "gzip")
	wrapped.ServeHTTP(rec3, req)

	if rec.Code != rec3.Code {
		t.Fatalf("Expected status code (%d) != actual (%d)", rec.Code, rec3.Code)
	}

	if rec.HeaderMap.Get("Content-Type") != rec3.HeaderMap.Get("Content-Type") {
		t.Fatal("Expected Content-Type header (%s) != actual (%s)", rec.HeaderMap.Get("Content-Type"), rec3.HeaderMap.Get("Content-Type"))
	}

	if rec3.HeaderMap.Get("Vary") != "Accept-Encoding" {
		t.Fatal("Missing Vary: Accept-Encoding header")
	}

	if rec3.HeaderMap.Get("Content-Encoding") != "gzip" {
		t.Fatal("Missing Content-Encoding: gzip header")
	}

	if rec.Body.Len() < rec3.Body.Len() {
		t.Fatalf("Payload actual size (%d) wasn't as small as expected (%d)", rec3.Body.Len(), rec.Body.Len())
	}

	gr, _ := gzip.NewReader(rec3.Body)
	defer gr.Close()

	data, _ := ioutil.ReadAll(gr)
	if string(data) != rec.Body.String() {
		t.Fatal("Expected body != actual when uncompressed")
	}
}
