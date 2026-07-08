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
	ListUsers(ctx context.Context, limit int, after, before uint) ([]model.User, error)
	ChangePassword(ctx context.Context, id uint, currentPassword, newPassword string) error
	// ChangeEmail requires the current password (email doubles as the login
	// credential) and rotates JwtID like ChangePassword, so a changed email
	// also logs out every existing session.
	ChangeEmail(ctx context.Context, id uint, currentPassword, newEmail string) (*model.User, error)
	UpdateRole(ctx context.Context, id uint, role model.Role) error
	DeleteUser(ctx context.Context, id uint) error
	SeedAdmin(ctx context.Context) error
	// SuspendComments sets or clears a user's commenting suspension. duration
	// is one of the presets ("1h", "1d", "1y"), a custom Go duration string
	// (e.g. "72h30m"), or "" / "none" to clear an existing suspension.
	SuspendComments(ctx context.Context, id uint, duration string) (*model.User, error)
	// UpdateProfile updates the caller's own display name and bio (self-service,
	// like ChangePassword — no permission check, only ever acts on the caller's
	// own account).
	UpdateProfile(ctx context.Context, id uint, name, bio string) (*model.User, error)
	SetAvatar(ctx context.Context, id uint, key, bucket string) (*model.User, error)
	// ClearAvatar unsets the caller's avatar, back to the initial-letter
	// placeholder. Leaves the old S3Object row/file in place, same as
	// replacing a series cover — the admin S3 cleanup tool catches it.
	ClearAvatar(ctx context.Context, id uint) (*model.User, error)
}

var commentSuspensionPresets = map[string]time.Duration{
	"1h": time.Hour,
	"1d": 24 * time.Hour,
	"1y": 365 * 24 * time.Hour,
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
		Role:         model.RoleReader,
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

func (s *userService) ListUsers(ctx context.Context, limit int, after, before uint) ([]model.User, error) {
	return s.repo.List(ctx, limit, after, before)
}

func (s *userService) ChangePassword(ctx context.Context, id uint, currentPassword, newPassword string) error {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashed)
	// Rotate JwtID so every existing session (including the one making this
	// request) is invalidated, matching the "password change logs everyone
	// out" behavior users expect.
	user.JwtID = uuid.New().String()

	return s.repo.Update(ctx, user)
}

func (s *userService) ChangeEmail(ctx context.Context, id uint, currentPassword, newEmail string) (*model.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return nil, errors.New("current password is incorrect")
	}

	if newEmail == user.Email {
		return user, nil
	}
	if existing, err := s.repo.GetByEmail(ctx, newEmail); err == nil && existing.ID != user.ID {
		return nil, errors.New("email is already in use")
	}

	user.Email = newEmail
	// Rotate JwtID so every existing session (including the one making this
	// request) is invalidated, matching ChangePassword's behavior — email
	// doubles as the login credential.
	user.JwtID = uuid.New().String()

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) UpdateRole(ctx context.Context, id uint, role model.Role) error {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	user.Role = role
	return s.repo.Update(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *userService) SuspendComments(ctx context.Context, id uint, duration string) (*model.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if duration == "" || duration == "none" {
		user.CommentSuspendedUntil = nil
	} else {
		d, ok := commentSuspensionPresets[duration]
		if !ok {
			d, err = time.ParseDuration(duration)
			if err != nil {
				return nil, errors.New("invalid duration: use \"1h\", \"1d\", \"1y\", a custom Go duration (e.g. \"72h30m\"), or \"none\" to clear")
			}
		}
		until := time.Now().Add(d)
		user.CommentSuspendedUntil = &until
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) UpdateProfile(ctx context.Context, id uint, name, bio string) (*model.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	user.Name = name
	user.Bio = bio
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) SetAvatar(ctx context.Context, id uint, key, bucket string) (*model.User, error) {
	return s.repo.SetAvatar(ctx, id, &model.S3Object{Key: key, Bucket: bucket})
}

func (s *userService) ClearAvatar(ctx context.Context, id uint) (*model.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	user.AvatarID = nil
	user.Avatar = nil
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
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
