package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/isa0-gh/reader/internal/model"
	"github.com/isa0-gh/reader/internal/service"
)

type ChapterHandler struct{ svc service.ChapterService }

func NewChapterHandler(svc service.ChapterService) *ChapterHandler {
	return &ChapterHandler{svc: svc}
}

func (h *ChapterHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	chapter, err := h.svc.GetChapter(r.Context(), uint(id))
	if err != nil {
		http.Error(w, "chapter not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chapter)
}

type createChapterRequest struct {
	SeriesID uint    `json:"series_id"`
	Number   float64 `json:"number"`
	Title    string  `json:"title"`
}

func (h *ChapterHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createChapterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.SeriesID == 0 {
		http.Error(w, "series_id is required", http.StatusBadRequest)
		return
	}

	chapter, err := h.svc.CreateChapter(r.Context(), &model.Chapter{
		SeriesID: req.SeriesID,
		Number:   req.Number,
		Title:    req.Title,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(chapter)
}
