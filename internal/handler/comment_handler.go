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

type CommentHandler struct {
	svc service.CommentService
}

func NewCommentHandler(svc service.CommentService) *CommentHandler {
	return &CommentHandler{svc: svc}
}

type createCommentRequest struct {
	Content  string `json:"content"`
	ParentID *uint  `json:"parent_id,omitempty"`
}

type updateCommentRequest struct {
	Content string `json:"content"`
}

type commentListResponse struct {
	Comments []model.Comment `json:"comments"`
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
}

// POST /api/v1/chapters/{chapterID}/comments
func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	var req createCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	comment, err := h.svc.CreateComment(r.Context(), user.ID, uint(chapterID), req.Content, req.ParentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

// GET /api/v1/chapters/{chapterID}/comments
func (h *CommentHandler) List(w http.ResponseWriter, r *http.Request) {
	chapterID, err := strconv.ParseUint(chi.URLParam(r, "chapterID"), 10, 32)
	if err != nil {
		http.Error(w, "invalid chapter id", http.StatusBadRequest)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	comments, total, err := h.svc.GetChapterComments(r.Context(), uint(chapterID), page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(commentListResponse{
		Comments: comments,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// GET /api/v1/comments/{commentID}/replies
func (h *CommentHandler) ListReplies(w http.ResponseWriter, r *http.Request) {
	commentID, err := strconv.ParseUint(chi.URLParam(r, "commentID"), 10, 32)
	if err != nil {
		http.Error(w, "invalid comment id", http.StatusBadRequest)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	replies, total, err := h.svc.GetCommentReplies(r.Context(), uint(commentID), page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(commentListResponse{
		Comments: replies,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// PATCH /api/v1/comments/{commentID}
func (h *CommentHandler) Update(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*model.User)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	commentID, err := strconv.ParseUint(chi.URLParam(r, "commentID"), 10, 32)
	if err != nil {
		http.Error(w, "invalid comment id", http.StatusBadRequest)
		return
	}

	var req updateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	comment, err := h.svc.UpdateComment(r.Context(), uint(commentID), user.ID, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}

// DELETE /api/v1/comments/{commentID}
func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*model.User)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	commentID, err := strconv.ParseUint(chi.URLParam(r, "commentID"), 10, 32)
	if err != nil {
		http.Error(w, "invalid comment id", http.StatusBadRequest)
		return
	}

	isAdmin := user.Role == "admin"
	if err := h.svc.DeleteComment(r.Context(), uint(commentID), user.ID, isAdmin); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// POST /api/v1/comments/{commentID}/like
func (h *CommentHandler) Like(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*model.User)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	commentID, err := strconv.ParseUint(chi.URLParam(r, "commentID"), 10, 32)
	if err != nil {
		http.Error(w, "invalid comment id", http.StatusBadRequest)
		return
	}

	if err := h.svc.LikeComment(r.Context(), uint(commentID), user.ID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DELETE /api/v1/comments/{commentID}/like
func (h *CommentHandler) Unlike(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*model.User)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	commentID, err := strconv.ParseUint(chi.URLParam(r, "commentID"), 10, 32)
	if err != nil {
		http.Error(w, "invalid comment id", http.StatusBadRequest)
		return
	}

	if err := h.svc.UnlikeComment(r.Context(), uint(commentID), user.ID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
