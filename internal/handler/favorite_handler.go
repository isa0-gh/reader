package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/isa0-gh/reader/internal/middleware"
	"github.com/isa0-gh/reader/internal/model"
	"github.com/isa0-gh/reader/internal/service"
)

type FavoriteHandler struct {
	svc service.FavoriteService
}

func NewFavoriteHandler(svc service.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{svc: svc}
}

// POST /api/v1/favorites/{seriesID}
func (h *FavoriteHandler) Add(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*model.User)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	seriesID, err := strconv.ParseUint(chi.URLParam(r, "seriesID"), 10, 32)
	if err != nil {
		http.Error(w, "invalid series id", http.StatusBadRequest)
		return
	}

	fav, err := h.svc.AddFavorite(r.Context(), user.ID, uint(seriesID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(fav)
}

// DELETE /api/v1/favorites/{seriesID}
func (h *FavoriteHandler) Remove(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*model.User)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	seriesID, err := strconv.ParseUint(chi.URLParam(r, "seriesID"), 10, 32)
	if err != nil {
		http.Error(w, "invalid series id", http.StatusBadRequest)
		return
	}

	if err := h.svc.RemoveFavorite(r.Context(), user.ID, uint(seriesID)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET /api/v1/favorites
func (h *FavoriteHandler) List(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*model.User)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	list, err := h.svc.GetUserFavorites(r.Context(), user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

// GET /api/v1/favorites/{seriesID}/check
func (h *FavoriteHandler) Check(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*model.User)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	seriesID, err := strconv.ParseUint(chi.URLParam(r, "seriesID"), 10, 32)
	if err != nil {
		http.Error(w, "invalid series id", http.StatusBadRequest)
		return
	}

	isFav, err := h.svc.IsFavorite(r.Context(), user.ID, uint(seriesID))
	if err != nil {
		http.Error(w, "check failed", http.StatusInternalServerError)
		return
	}

	count, _ := h.svc.GetFavoriteCount(r.Context(), uint(seriesID))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"is_favorite":     isFav,
		"total_favorites": count,
	})
}
