package api

import (
	"log/slog"
	"net/http"
)

// HandleWebSocket is retained for backwards compatibility. The browser app uses
// the HTTP AI endpoint and online room polling, so no external WebSocket package
// is required for deployment.
func HandleWebSocket(logger *slog.Logger, w http.ResponseWriter, r *http.Request) {
	logger.Info("legacy websocket endpoint requested", "path", r.URL.Path)
	writeJSON(w, http.StatusNotImplemented, ErrorResponse{
		Error:   "not_implemented",
		Message: "UCI WebSocket is disabled in the web deployment build",
	})
}
