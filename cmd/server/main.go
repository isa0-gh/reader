package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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
	if err := db.AutoMigrate(
		&model.User{},
		&model.Series{},
		&model.Chapter{},
		&model.S3Object{},
	); err != nil {
		log.Fatalf("could not migrate database: %v", err)
	}

	// Initialize Repository, Service, and Handler
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	seriesRepo := repository.NewSeriesRepository(db)
	seriesSvc := service.NewSeriesService(seriesRepo)
	seriesHandler := handler.NewSeriesHandler(seriesSvc)

	chapterRepo := repository.NewChapterRepository(db)
	chapterSvc := service.NewChapterService(chapterRepo)
	chapterHandler := handler.NewChapterHandler(chapterSvc)

	// Setup Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Basic CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"}, // Adjust for production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// API Routes
	r.Route("/api/v1", func(r chi.Router) {
		// Auth
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", userHandler.Register)
			r.Post("/login", userHandler.Login)
		})

		// Users (protected)
		r.Route("/users", func(r chi.Router) {
			r.Use(appMiddleware.JWTMiddleware(userRepo))
			r.Get("/", userHandler.List)
			r.Get("/{id}", userHandler.Get)
		})

		// Series
		r.Get("/series/{id}", seriesHandler.Get)

		// Chapters
		r.Get("/chapters/{id}", chapterHandler.Get)
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
