package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-shrotener/internal/config"
	handlers "url-shrotener/internal/handlers/http"
	"url-shrotener/internal/storage"
	"url-shrotener/tools"
)

func main() {
	cfg := config.Load()

	log := slog.New(
		// LevelDebug for development environment only
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	log.Info("starting url shortener")
	log.Debug("debug messages are enabled")

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	switch cfg.HTTPServer.StorageType {
	case "db":
		strg, err := storage.NewDB(cfg)
		if err != nil {
			log.Error("failed to init database", tools.LogAttr("error", err.Error()))
			os.Exit(1)
		}
		defer strg.DB.Close()

		log.Info(
			"successfully connected to database",
			slog.String("host", cfg.DB.Host),
			slog.String("port", cfg.DB.Port),
		)

		router.Route("/api", func(r chi.Router) {
			r.Post("/http", handlers.Save(log, strg, cfg.URLLength, cfg.URLAlphabet))
			r.Get("/http", handlers.Get(log, strg))
		})

	case "inmem":
		strg := storage.NewInMemStorage()
		log.Info("successfully created inmemory storage")

		router.Route("/api", func(r chi.Router) {
			r.Post("/http", handlers.Save(log, strg, cfg.URLLength, cfg.URLAlphabet))
			r.Get("/http", handlers.Get(log, strg))
		})
	default:
		log.Error("storage type isn't specified: --storage=db|inmem")
		os.Exit(1)
	}

	srv := &http.Server{
		Addr:    cfg.HTTPServer.Host + ":" + cfg.HTTPServer.Port,
		Handler: router,
	}

	stopped := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Error("HTTP Server Shutdown Error ", slog.String("error", err.Error()))
		}
		close(stopped)
	}()

	log.Info(
		"Starting HTTP server",
		slog.String("host", cfg.HTTPServer.Host),
		slog.String("port", cfg.HTTPServer.Port),
	)

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Error("HTTP server ListenAndServe error ", tools.LogAttr("error", err.Error()))
	}

	<-stopped

	log.Info("Server has been gracefully stopped")
}
