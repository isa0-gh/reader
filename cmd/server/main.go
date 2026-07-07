package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/isa0-gh/reader/internal/config"
	"github.com/isa0-gh/reader/internal/database"
	"github.com/isa0-gh/reader/internal/handler"
	appMiddleware "github.com/isa0-gh/reader/internal/middleware"
	"github.com/isa0-gh/reader/internal/model"
	"github.com/isa0-gh/reader/internal/repository"
	"github.com/isa0-gh/reader/internal/service"
	"github.com/isa0-gh/reader/internal/storage"
)

// clientIPKey rate-limits by the client IP resolved by
// middleware.ClientIPFromXFFTrustedProxies, rather than trusting
// X-Forwarded-For/RemoteAddr directly (both are spoofable without it).
func clientIPKey(r *http.Request) (string, error) {
	ip := middleware.GetClientIP(r.Context())
	if ip == "" {
		return "", errors.New("client ip not resolved")
	}
	return httprate.CanonicalizeIP(ip), nil
}

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	for _, m := range []any{&model.User{}, &model.Series{}, &model.Chapter{}, &model.S3Object{}} {
		if err := db.AutoMigrate(m); err != nil {
			log.Fatalf("could not migrate database: %v", err)
		}
	}

	// Initialize Repository, Service, and Handler
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc, cfg)

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
	// The backend is only ever reached through the frontend nginx container
	// (docker-compose publishes no backend port), so there is exactly one
	// trusted hop between us and the client: resolve the real client IP from
	// the outermost X-Forwarded-For entry that proxy adds.
	r.Use(middleware.ClientIPFromXFFTrustedProxies(1))

	// Basic CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"}, // Adjust for production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Maintenance mode: block all API traffic except health checks and the
	// config endpoint (frontend needs it to render the maintenance notice).
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cfg.Maintenance && r.URL.Path != "/health" && r.URL.Path != "/api/v1/config" {
				http.Error(w, "service is under maintenance", http.StatusServiceUnavailable)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	// API Routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/config", configHandler.Get)

		// Auth (rate limited per-IP to deter brute-force login/register),
		// keyed off the client IP resolved by ClientIPFromXFFTrustedProxies.
		r.Route("/auth", func(r chi.Router) {
			r.Use(httprate.LimitBy(10, time.Minute, clientIPKey))
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
			r.With(appMiddleware.RequirePermission("series:update")).Put("/series/{id}", seriesHandler.Update)
			r.With(appMiddleware.RequirePermission("series:delete")).Delete("/series/{id}", seriesHandler.Delete)
		})

		// Chapters
		r.Route("/chapters", func(r chi.Router) {
			r.Get("/{id}", chapterHandler.Get)
			r.Group(func(r chi.Router) {
				r.Use(appMiddleware.JWTMiddleware(userRepo))
				r.With(appMiddleware.RequirePermission("chapter:create")).Post("/", chapterHandler.Create)
				r.With(appMiddleware.RequirePermission("chapter:create")).Post("/{id}/pages", chapterHandler.UploadPages)
				// Update/ReorderPages/DeletePage/Delete: chapterHandler.authorize
				// enforces chapter:X or chapter:X:own (matching uploader) per-request,
				// since ownership can only be resolved after loading the chapter.
				r.Put("/{id}", chapterHandler.Update)
				r.Put("/{id}/pages", chapterHandler.ReorderPages)
				r.Delete("/{id}/pages/{pageId}", chapterHandler.DeletePage)
				r.Delete("/{id}", chapterHandler.Delete)
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
