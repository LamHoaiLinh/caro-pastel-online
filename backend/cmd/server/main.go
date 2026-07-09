package main

import (
	"caro-ai-pvp/internal/api"
	"caro-ai-pvp/internal/domain"
	"caro-ai-pvp/internal/persistence"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"syscall"
	"time"
)

func main() {
	debug.SetMemoryLimit(domain.HeapHardLimitBytes)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	store := api.NewInMemoryStore()

	dbPath := filepath.Join(".", "data", "matches.json")
	if v := os.Getenv("MATCH_DB_PATH"); v != "" {
		dbPath = v
	}
	matchStore, err := persistence.NewMatchStore(dbPath)
	if err != nil {
		logger.Error("failed to open match database", "err", err, "path", dbPath)
		os.Exit(1)
	}

	handler := api.NewHandler(store, matchStore, logger)
	server := api.NewServer(handler, logger)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5207"
	}

	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: server,
	}

	serverErr := make(chan error, 1)
	go func() {
		logger.Info("server starting", "addr", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	cleanupTicker := time.NewTicker(5 * time.Minute)
	go func() {
		for range cleanupTicker.C {
			removed := store.CleanupCompleted()
			onlineRemoved := handler.CleanupOnlineRooms()
			if removed > 0 || onlineRemoved > 0 {
				logger.Info("cleanup", "gamesRemoved", removed, "onlineRoomsRemoved", onlineRemoved)
			}
		}
	}()

	select {
	case err := <-serverErr:
		logger.Error("server failed to start", "err", err)
		matchStore.Close()
		fmt.Fprintf(os.Stderr, "Fatal: %v\n", err)
		os.Exit(1)
	case <-quit:
	}

	logger.Info("shutting down")
	cleanupTicker.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("shutdown error", "err", err)
	}

	remaining := store.CleanupAll()
	onlineRemaining := handler.CleanupAllOnlineRooms()
	if remaining > 0 || onlineRemaining > 0 {
		logger.Info("shutdown cleanup", "remaining", remaining, "onlineRemaining", onlineRemaining)
	}
	matchStore.Close()

	fmt.Println("Server stopped")
}
