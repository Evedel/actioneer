package healthz

import (
	"log/slog"
	"net/http"
)

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Debug("/healthz request")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
