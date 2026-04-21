package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/isa0-gh/reader/internal/config"
	"github.com/isa0-gh/reader/internal/database"
	"github.com/isa0-gh/reader/internal/handler"
	appMiddleware "github.com/isa0-gh/reader/internal/middleware"
	"github.com/isa0-gh/reader/internal/model"
	"github.com/isa0-gh/reader/internal/repository"
	"github.com/isa0-gh/reader/internal/service"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	// Auto-migrate models
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("could not migrate database: %v", err)
	}

	// Initialize Repository, Service, and Handler
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	// Setup Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// API Routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.Register)

			// Protected routes
			r.Group(func(r chi.Router) {
				r.Use(appMiddleware.JWTMiddleware(userRepo))
				r.Get("/", userHandler.List)
				r.Get("/{id}", userHandler.Get)
			})
		})
	})

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	log.Printf("Server starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
