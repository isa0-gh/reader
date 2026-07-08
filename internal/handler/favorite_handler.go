package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/isa0-gh/reader/internal/httpx"
	appmw "github.com/isa0-gh/reader/internal/middleware"
	"github.com/isa0-gh/reader/internal/model"
	"github.com/isa0-gh/reader/internal/service"
)

type FavoriteHandler struct{ svc service.FavoriteService }

func NewFavoriteHandler(svc service.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{svc: svc}
}

func currentUser(r *http.Request) (*model.User, bool) {
	user, ok := r.Context().Value(appmw.UserContextKey).(*model.User)
	return user, ok
}

// List returns the authenticated user's favorited series, each with its
// series (cover included) and last-read-chapter progress preloaded.
func (h *FavoriteHandler) List(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(r)
	if !ok {
		httpx.WriteError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	limit := 50
	offset := 0
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}

	list, total, err := h.svc.ListFavorites(r.Context(), user.ID, limit, offset)
	if err != nil {
		httpx.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"items": list,
		"total": total,
	})
}

type addFavoriteRequest struct {
	SeriesID uint `json:"series_id"`
}

func (h *FavoriteHandler) Add(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(r)
	if !ok {
		httpx.WriteError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req addFavoriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.SeriesID == 0 {
		httpx.WriteError(w, "series_id is required", http.StatusBadRequest)
		return
	}

	fav, err := h.svc.AddFavorite(r.Context(), user.ID, req.SeriesID)
	if err != nil {
		httpx.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(fav)
}

func (h *FavoriteHandler) Remove(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(r)
	if !ok {
		httpx.WriteError(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	seriesID, err := strconv.ParseUint(chi.URLParam(r, "seriesId"), 10, 32)
	if err != nil {
		httpx.WriteError(w, "invalid series id", http.StatusBadRequest)
		return
	}
	if err := h.svc.RemoveFavorite(r.Context(), user.ID, uint(seriesID)); err != nil {
		httpx.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type updateProgressRequest struct {
	ChapterID uint `json:"chapter_id"`
}

// UpdateProgress is a best-effort ping: the frontend calls it on every
// chapter view regardless of favorite status, so it silently no-ops (204,
// no error) when the series isn't favorited rather than requiring the
// caller to check first.
func (h *FavoriteHandler) UpdateProgress(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(r)
	if !ok {
		httpx.WriteError(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	seriesID, err := strconv.ParseUint(chi.URLParam(r, "seriesId"), 10, 32)
	if err != nil {
		httpx.WriteError(w, "invalid series id", http.StatusBadRequest)
		return
	}

	var req updateProgressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.ChapterID == 0 {
		httpx.WriteError(w, "chapter_id is required", http.StatusBadRequest)
		return
	}

	if err := h.svc.UpdateProgress(r.Context(), user.ID, uint(seriesID), req.ChapterID); err != nil {
		httpx.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
