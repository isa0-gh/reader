package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/isa0-gh/reader/internal/config"
	"github.com/isa0-gh/reader/internal/httpx"
	appmw "github.com/isa0-gh/reader/internal/middleware"
	"github.com/isa0-gh/reader/internal/model"
	"github.com/isa0-gh/reader/internal/service"
)

type UserHandler struct {
	svc service.UserService
	cfg *config.Config
}

func NewUserHandler(svc service.UserService, cfg *config.Config) *UserHandler {
	return &UserHandler{svc: svc, cfg: cfg}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if h.cfg.RegisterDisabled {
		httpx.WriteError(w, "registration is disabled", http.StatusForbidden)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.svc.RegisterUser(r.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		httpx.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  any    `json:"user"`
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if h.cfg.LoginDisabled {
		httpx.WriteError(w, "login is disabled", http.StatusForbidden)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, user, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		httpx.WriteError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Token: token,
		User:  user,
	})
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		httpx.WriteError(w, "invalid user id", http.StatusBadRequest)
		return
	}

	user, err := h.svc.GetUser(r.Context(), uint(id))
	if err != nil {
		httpx.WriteError(w, "user not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := 50
	var after, before uint

	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}
	if v := r.URL.Query().Get("after"); v != "" {
		if n, err := strconv.ParseUint(v, 10, 32); err == nil {
			after = uint(n)
		}
	}
	if v := r.URL.Query().Get("before"); v != "" {
		if n, err := strconv.ParseUint(v, 10, 32); err == nil {
			before = uint(n)
		}
	}

	users, err := h.svc.ListUsers(r.Context(), limit, after, before)
	if err != nil {
		httpx.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

type changePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// PATCH /api/v1/users/me/password — any authenticated user changes their own
// password. Unlike UpdateRole/Delete this isn't admin-gated: it only ever
// acts on the caller's own account, taken from the JWT, not a URL param.
func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(appmw.UserContextKey).(*model.User)
	if !ok {
		httpx.WriteError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req changePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.NewPassword == "" {
		httpx.WriteError(w, "new_password is required", http.StatusBadRequest)
		return
	}

	if err := h.svc.ChangePassword(r.Context(), user.ID, req.CurrentPassword, req.NewPassword); err != nil {
		httpx.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		httpx.WriteError(w, "invalid user id", http.StatusBadRequest)
		return
	}
	var body struct {
		Role model.Role `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpx.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.svc.UpdateRole(r.Context(), uint(id), body.Role); err != nil {
		httpx.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type suspendCommentsRequest struct {
	Duration string `json:"duration"`
}

// SuspendComments is admin/mod-only (comment:suspend). duration accepts the
// presets "1h"/"1d"/"1y", a custom Go duration string (e.g. "72h30m"), or ""
// / "none" to lift an existing suspension.
func (h *UserHandler) SuspendComments(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		httpx.WriteError(w, "invalid user id", http.StatusBadRequest)
		return
	}
	var req suspendCommentsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := h.svc.SuspendComments(r.Context(), uint(id), req.Duration)
	if err != nil {
		httpx.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		httpx.WriteError(w, "invalid user id", http.StatusBadRequest)
		return
	}
	if err := h.svc.DeleteUser(r.Context(), uint(id)); err != nil {
		httpx.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
