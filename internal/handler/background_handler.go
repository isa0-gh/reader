package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/isa0-gh/reader/internal/httpx"
	"github.com/isa0-gh/reader/internal/service"
)

type BackgroundHandler struct{ svc service.BackgroundService }

func NewBackgroundHandler(svc service.BackgroundService) *BackgroundHandler {
	return &BackgroundHandler{svc: svc}
}

// List is public (no auth) — every visitor's page needs the wallpaper set to
// render the site background, not just admins.
func (h *BackgroundHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.ListBackgrounds(r.Context())
	if err != nil {
		httpx.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

type addBackgroundRequest struct {
	Key    string `json:"key"`
	Bucket string `json:"bucket"`
}

func (h *BackgroundHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req addBackgroundRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Key == "" || req.Bucket == "" {
		httpx.WriteError(w, "key and bucket are required", http.StatusBadRequest)
		return
	}

	bg, err := h.svc.AddBackground(r.Context(), req.Key, req.Bucket)
	if err != nil {
		httpx.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bg)
}

func (h *BackgroundHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		httpx.WriteError(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.svc.DeleteBackground(r.Context(), uint(id)); err != nil {
		httpx.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
