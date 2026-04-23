package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/isa0-gh/reader/internal/model"
	"github.com/isa0-gh/reader/internal/service"
)

type SeriesHandler struct{ svc service.SeriesService }

func NewSeriesHandler(svc service.SeriesService) *SeriesHandler {
	return &SeriesHandler{svc: svc}
}

func (h *SeriesHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.ListSeries(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *SeriesHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	series, err := h.svc.GetSeries(r.Context(), uint(id))
	if err != nil {
		http.Error(w, "series not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(series)
}

type createSeriesRequest struct {
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	CoverImage  string `json:"cover_image"`
	Author      string `json:"author"`
	Artist      string `json:"artist"`
	Status      string `json:"status"`
}

func (h *SeriesHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createSeriesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Title == "" || req.Slug == "" {
		http.Error(w, "title and slug are required", http.StatusBadRequest)
		return
	}

	status := model.SeriesStatus(req.Status)
	if status == "" {
		status = model.StatusOngoing
	}

	series, err := h.svc.CreateSeries(r.Context(), &model.Series{
		Title:       req.Title,
		Slug:        req.Slug,
		Description: req.Description,
		CoverImage:  req.CoverImage,
		Author:      req.Author,
		Artist:      req.Artist,
		Status:      status,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(series)
}
