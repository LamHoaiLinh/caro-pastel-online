package api

import (
	"log/slog"
	"net/http"
)

func NewServer(handler *Handler, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{"status": "ok"})
	})

	mux.HandleFunc("POST /api/game/new", handler.CreateGame)
	mux.HandleFunc("GET /api/game/{id}", handler.GetGame)
	mux.HandleFunc("POST /api/game/{id}/move", handler.MakeMove)
	mux.HandleFunc("POST /api/game/{id}/ai-move", handler.MakeAIMove)
	mux.HandleFunc("POST /api/game/{id}/undo", handler.UndoMove)
	mux.HandleFunc("DELETE /api/game/{id}", handler.DeleteGame)
	mux.HandleFunc("GET /ws/uci", func(w http.ResponseWriter, r *http.Request) {
		HandleWebSocket(logger, w, r)
	})
	mux.HandleFunc("POST /api/online/create", handler.CreateOnlineRoom)
	mux.HandleFunc("POST /api/online/{code}/join", handler.JoinOnlineRoom)
	mux.HandleFunc("GET /api/online/{code}", handler.GetOnlineRoom)
	mux.HandleFunc("POST /api/online/{code}/move", handler.MakeOnlineMove)

	var h http.Handler = mux
	h = CORSMiddleware(h)
	h = LoggingMiddleware(logger, h)
	h = RecoveryMiddleware(logger, h)

	return h
}
