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

type ChapterHandler struct{ svc service.ChapterService }

func NewChapterHandler(svc service.ChapterService) *ChapterHandler {
	return &ChapterHandler{svc: svc}
}

// authorize checks that the caller either holds the blanket "chapter:<action>"
// permission or holds "chapter:<action>:own" and uploaded the chapter themselves.
func (h *ChapterHandler) authorize(w http.ResponseWriter, r *http.Request, id uint, action string) (*model.Chapter, bool) {
	user, ok := r.Context().Value(appmw.UserContextKey).(*model.User)
	if !ok {
		httpx.WriteError(w, "unauthorized", http.StatusUnauthorized)
		return nil, false
	}

	full := "chapter:" + action
	if !user.HasPermission(full) && !user.HasPermission(full+":own") {
		httpx.WriteError(w, "forbidden", http.StatusForbidden)
		return nil, false
	}

	chapter, err := h.svc.GetChapter(r.Context(), id)
	if err != nil {
		httpx.WriteError(w, "chapter not found", http.StatusNotFound)
		return nil, false
	}

	if user.HasPermission(full) || (chapter.UploaderID != nil && *chapter.UploaderID == user.ID) {
		return chapter, true
	}

	httpx.WriteError(w, "forbidden", http.StatusForbidden)
	return nil, false
}

func (h *ChapterHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		httpx.WriteError(w, "invalid id", http.StatusBadRequest)
		return
	}

	chapter, err := h.svc.GetChapter(r.Context(), uint(id))
	if err != nil {
		httpx.WriteError(w, "chapter not found", http.StatusNotFound)
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
		httpx.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.SeriesID == 0 {
		httpx.WriteError(w, "series_id is required", http.StatusBadRequest)
		return
	}

	var uploaderID *uint
	if user, ok := r.Context().Value(appmw.UserContextKey).(*model.User); ok {
		uploaderID = &user.ID
	}

	chapter, err := h.svc.CreateChapter(r.Context(), &model.Chapter{
		SeriesID:   req.SeriesID,
		UploaderID: uploaderID,
		Number:     req.Number,
		Title:      req.Title,
	})
	if err != nil {
		httpx.WriteError(w, err.Error(), http.StatusInternalServerError)
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
		httpx.WriteError(w, "invalid id", http.StatusBadRequest)
		return
	}
	if _, ok := h.authorize(w, r, uint(id), "update"); !ok {
		return
	}

	var req updateChapterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	chapter, err := h.svc.UpdateChapter(r.Context(), &model.Chapter{
		ID:     uint(id),
		Number: req.Number,
		Title:  req.Title,
	})
	if err != nil {
		httpx.WriteError(w, err.Error(), http.StatusInternalServerError)
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
		httpx.WriteError(w, "invalid id", http.StatusBadRequest)
		return
	}
	if _, ok := h.authorize(w, r, uint(id), "update"); !ok {
		return
	}

	var req reorderPagesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	numbers := make(map[uint]int, len(req.Pages))
	for _, p := range req.Pages {
		numbers[p.ID] = p.PageNumber
	}

	if err := h.svc.UpdatePageNumbers(r.Context(), uint(id), numbers); err != nil {
		httpx.WriteError(w, err.Error(), http.StatusInternalServerError)
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
		httpx.WriteError(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req uploadPagesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, err.Error(), http.StatusBadRequest)
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
		httpx.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *ChapterHandler) DeletePage(w http.ResponseWriter, r *http.Request) {
	chapterID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		httpx.WriteError(w, "invalid chapter id", http.StatusBadRequest)
		return
	}
	pageID, err := strconv.ParseUint(chi.URLParam(r, "pageId"), 10, 32)
	if err != nil {
		httpx.WriteError(w, "invalid page id", http.StatusBadRequest)
		return
	}
	if _, ok := h.authorize(w, r, uint(chapterID), "update"); !ok {
		return
	}
	if err := h.svc.DeletePage(r.Context(), uint(chapterID), uint(pageID)); err != nil {
		httpx.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ChapterHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		httpx.WriteError(w, "invalid id", http.StatusBadRequest)
		return
	}
	if _, ok := h.authorize(w, r, uint(id), "delete"); !ok {
		return
	}
	if err := h.svc.DeleteChapter(r.Context(), uint(id)); err != nil {
		httpx.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
