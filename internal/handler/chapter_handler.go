package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
