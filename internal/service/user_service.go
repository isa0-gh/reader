package service

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/isa0-gh/reader/internal/middleware"
	"github.com/isa0-gh/reader/internal/model"
	"github.com/isa0-gh/reader/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUser(ctx context.Context, email, password, name string) (*model.User, error)
	Login(ctx context.Context, email, password string) (string, *model.User, error)
	GetUser(ctx context.Context, id uint) (*model.User, error)
	ListUsers(ctx context.Context) ([]model.User, error)
	SeedAdmin(ctx context.Context) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) RegisterUser(ctx context.Context, email, password, name string) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		Name:         name,
		JwtID:        uuid.New().String(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(ctx context.Context, email, password string) (string, *model.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	// Update JwtID on login for revocation logic
	user.JwtID = uuid.New().String()
	if err := s.repo.Update(ctx, user); err != nil {
		return "", nil, err
	}

	// Generate JWT
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default_secret_change_me"
	}

	claims := &middleware.Claims{
		UserID: user.ID,
		JwtID:  user.JwtID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", nil, err
	}

	return tokenString, user, nil
}

func (s *userService) GetUser(ctx context.Context, id uint) (*model.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) ListUsers(ctx context.Context) ([]model.User, error) {
	return s.repo.List(ctx)
}

func (s *userService) SeedAdmin(ctx context.Context) error {
	count, err := s.repo.Count(ctx)
	if err != nil || count > 0 {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.Create(ctx, &model.User{
		Email:        "admin@localhost",
		PasswordHash: string(hash),
		Name:         "Admin",
		Role:         model.RoleAdmin,
		JwtID:        uuid.New().String(),
	})
}
