package healthz

import (
	th "actioneer/internal/testing_helper"
	"net/http"
	"testing"
)

type fakeResponseWriter struct {
	header http.Header
	code   int
	body   []byte
}

func (w *fakeResponseWriter) Header() http.Header {
	return w.header
}
func (w *fakeResponseWriter) Write(b []byte) (int, error) {
	w.body = b
	return len(b), nil
}
func (w *fakeResponseWriter) WriteHeader(code int) {
	w.code = code
}

func TestServeHTTP(t *testing.T) {
	w := &fakeResponseWriter{}
	r := &http.Request{}
	ServeHTTP(w, r)

	th.AssertEqual(t, w.code, http.StatusOK)
	th.AssertEqual(t, string(w.body), "ok")
}
