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

type ProgressHandler struct{ svc service.ProgressService }

func NewProgressHandler(svc service.ProgressService) *ProgressHandler {
	return &ProgressHandler{svc: svc}
}

type updateProgressRequest struct {
	PageIndex int  `json:"page_index"`
	Completed bool `json:"completed"`
}

// PUT /api/v1/progress/chapters/{chapterID}
func (h *ProgressHandler) Update(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*model.User)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	chapterID, err := strconv.ParseUint(chi.URLParam(r, "chapterID"), 10, 32)
	if err != nil {
		http.Error(w, "invalid chapter id", http.StatusBadRequest)
		return
	}

	var req updateProgressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	progress, err := h.svc.UpdateProgress(r.Context(), user.ID, uint(chapterID), req.PageIndex, req.Completed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(progress)
}

// GET /api/v1/progress/chapters/{chapterID}
func (h *ProgressHandler) Get(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*model.User)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	chapterID, err := strconv.ParseUint(chi.URLParam(r, "chapterID"), 10, 32)
	if err != nil {
		http.Error(w, "invalid chapter id", http.StatusBadRequest)
		return
	}

	progress, err := h.svc.GetProgress(r.Context(), user.ID, uint(chapterID))
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(progress)
}

// GET /api/v1/progress/me
func (h *ProgressHandler) ListMine(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*model.User)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	list, err := h.svc.GetUserProgress(r.Context(), user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

// GET /api/v1/progress/series/{seriesID}
func (h *ProgressHandler) GetSeriesProgress(w http.ResponseWriter, r *http.Request) {
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

	list, err := h.svc.GetSeriesProgress(r.Context(), user.ID, uint(seriesID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}
