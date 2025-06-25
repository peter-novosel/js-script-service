package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/peter-novosel/js-script-service/internal/config"
	"github.com/peter-novosel/js-script-service/internal/logger"
	"github.com/peter-novosel/js-script-service/internal/executor"
	"github.com/peter-novosel/js-script-service/internal/admin"
)

func Setup(cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	}))

	// Health check
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Dynamic script execution
	r.Route("/scripts", func(r chi.Router) {
		r.Get("/{slug}", executor.HandleExecuteScript)
		r.Post("/{slug}", executor.HandleExecuteScript)
	})

	// Admin routes
	r.Route("/admin", func(r chi.Router) {
		r.Post("/scripts", admin.CreateOrUpdateScript)
		r.Get("/scripts", admin.ListScripts)
	})

	

	logger.Init().Info("Router setup complete")
	return r
}
