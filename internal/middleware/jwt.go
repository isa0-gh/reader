package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/isa0-gh/reader/internal/repository"
)

type contextKey string

const UserContextKey contextKey = "user"

type Claims struct {
	UserID uint   `json:"user_id"`
	JwtID  string `json:"jti"`
	jwt.RegisteredClaims
}

func JWTMiddleware(userRepo repository.UserRepository) func(http.Handler) http.Handler {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default_secret_change_me"
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			tokenString := bearerToken[1]
			claims := &Claims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Revocation logic: Check if JwtID matches the one in the database
			user, err := userRepo.GetByID(r.Context(), claims.UserID)
			if err != nil || user.JwtID != claims.JwtID {
				http.Error(w, "Token has been revoked", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
