package main

import (
	"context"
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
	"github.com/isa0-gh/reader/internal/storage"
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
		&model.S3Object{},
		&model.Series{},
		&model.Chapter{},
		&model.User{},
	); err != nil {
		log.Fatalf("could not migrate database: %v", err)
	}

	// Initialize Repository, Service, and Handler
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	// Seed first admin if no users exist
	if err := userSvc.SeedAdmin(context.Background()); err != nil {
		log.Fatalf("could not seed admin: %v", err)
	}

	seriesRepo := repository.NewSeriesRepository(db)
	seriesSvc := service.NewSeriesService(seriesRepo)
	seriesHandler := handler.NewSeriesHandler(seriesSvc)

	chapterRepo := repository.NewChapterRepository(db)
	chapterSvc := service.NewChapterService(chapterRepo)
	chapterHandler := handler.NewChapterHandler(chapterSvc)

	s3Client := storage.NewS3Client(cfg)
	uploadHandler := handler.NewUploadHandler(s3Client)
	configHandler := handler.NewConfigHandler(cfg)
	s3CleanHandler := handler.NewS3CleanHandler(db, s3Client)

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
		r.Get("/config", configHandler.Get)

		// Auth
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", userHandler.Register)
			r.Post("/login", userHandler.Login)
		})

		// Users (admin only)
		r.Route("/users", func(r chi.Router) {
			r.Use(appMiddleware.JWTMiddleware(userRepo))
			r.Use(appMiddleware.RequirePermission("user:list"))
			r.Get("/", userHandler.List)
			r.Get("/{id}", userHandler.Get)
			r.With(appMiddleware.RequirePermission("user:update")).Patch("/{id}/role", userHandler.UpdateRole)
			r.With(appMiddleware.RequirePermission("user:delete")).Delete("/{id}", userHandler.Delete)
		})

		// Admin: S3 cleanup (admin only)
		r.Route("/admin/s3", func(r chi.Router) {
			r.Use(appMiddleware.JWTMiddleware(userRepo))
			r.Use(appMiddleware.RequirePermission("user:list"))
			r.Get("/orphaned", s3CleanHandler.List)
			r.Delete("/orphaned", s3CleanHandler.Purge)
		})

		// Upload (presign) — requires auth
		r.Group(func(r chi.Router) {
			r.Use(appMiddleware.JWTMiddleware(userRepo))
			r.Post("/upload/presign", uploadHandler.Presign)
		})

		// Series
		r.Get("/series", seriesHandler.List)
		r.Get("/series/{id}", seriesHandler.Get)
		r.Group(func(r chi.Router) {
			r.Use(appMiddleware.JWTMiddleware(userRepo))
			r.With(appMiddleware.RequirePermission("series:create")).Post("/series", seriesHandler.Create)
			r.With(appMiddleware.RequirePermission("series:delete")).Delete("/series/{id}", seriesHandler.Delete)
		})

		// Chapters
		r.Route("/chapters", func(r chi.Router) {
			r.Get("/{id}", chapterHandler.Get)
			r.Group(func(r chi.Router) {
				r.Use(appMiddleware.JWTMiddleware(userRepo))
				r.With(appMiddleware.RequirePermission("chapter:create")).Post("/", chapterHandler.Create)
				r.With(appMiddleware.RequirePermission("chapter:create")).Post("/{id}/pages", chapterHandler.UploadPages)
				r.With(appMiddleware.RequirePermission("chapter:update")).Delete("/{id}/pages/{pageId}", chapterHandler.DeletePage)
				r.With(appMiddleware.RequirePermission("chapter:delete")).Delete("/{id}", chapterHandler.Delete)
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
