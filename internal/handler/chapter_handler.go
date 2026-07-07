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

type updateChapterRequest struct {
	Number float64 `json:"number"`
	Title  string  `json:"title"`
}

func (h *ChapterHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req updateChapterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	chapter, err := h.svc.UpdateChapter(r.Context(), &model.Chapter{
		ID:     uint(id),
		Number: req.Number,
		Title:  req.Title,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chapter)
}

type reorderPagesRequest struct {
	Pages []struct {
		ID         uint `json:"id"`
		PageNumber int  `json:"page_number"`
	} `json:"pages"`
}

func (h *ChapterHandler) ReorderPages(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req reorderPagesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	numbers := make(map[uint]int, len(req.Pages))
	for _, p := range req.Pages {
		numbers[p.ID] = p.PageNumber
	}

	if err := h.svc.UpdatePageNumbers(r.Context(), uint(id), numbers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type uploadPagesRequest struct {
	Pages []struct {
		Key        string `json:"key"`
		Bucket     string `json:"bucket"`
		PageNumber int    `json:"page_number"`
	} `json:"pages"`
}

func (h *ChapterHandler) UploadPages(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req uploadPagesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pages := make([]model.S3Object, len(req.Pages))
	for i, p := range req.Pages {
		pages[i] = model.S3Object{
			Key:        p.Key,
			Bucket:     p.Bucket,
			PageNumber: p.PageNumber,
		}
	}

	if err := h.svc.AddPages(r.Context(), uint(id), pages); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *ChapterHandler) DeletePage(w http.ResponseWriter, r *http.Request) {
	chapterID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "invalid chapter id", http.StatusBadRequest)
		return
	}
	pageID, err := strconv.ParseUint(chi.URLParam(r, "pageId"), 10, 32)
	if err != nil {
		http.Error(w, "invalid page id", http.StatusBadRequest)
		return
	}
	if err := h.svc.DeletePage(r.Context(), uint(chapterID), uint(pageID)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ChapterHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.svc.DeleteChapter(r.Context(), uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
