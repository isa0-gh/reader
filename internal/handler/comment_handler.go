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

type CommentHandler struct{ svc service.CommentService }

func NewCommentHandler(svc service.CommentService) *CommentHandler {
	return &CommentHandler{svc: svc}
}

func (h *CommentHandler) List(w http.ResponseWriter, r *http.Request) {
	chapterID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		httpx.WriteError(w, "invalid chapter id", http.StatusBadRequest)
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

	list, total, err := h.svc.ListComments(r.Context(), uint(chapterID), limit, offset)
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

type createCommentRequest struct {
	Body string `json:"body"`
}

func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	chapterID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		httpx.WriteError(w, "invalid chapter id", http.StatusBadRequest)
		return
	}
	user, ok := r.Context().Value(appmw.UserContextKey).(*model.User)
	if !ok {
		httpx.WriteError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req createCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	comment, err := h.svc.CreateComment(r.Context(), uint(chapterID), user, req.Body)
	if err != nil {
		httpx.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

// Delete allows the comment's own author, or anyone holding the blanket
// "comment:delete" permission (moderator/admin), to remove a comment.
func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		httpx.WriteError(w, "invalid id", http.StatusBadRequest)
		return
	}
	user, ok := r.Context().Value(appmw.UserContextKey).(*model.User)
	if !ok {
		httpx.WriteError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	comment, err := h.svc.GetComment(r.Context(), uint(id))
	if err != nil {
		httpx.WriteError(w, "comment not found", http.StatusNotFound)
		return
	}

	if comment.UserID != user.ID && !user.HasPermission("comment:delete") {
		httpx.WriteError(w, "forbidden", http.StatusForbidden)
		return
	}

	if err := h.svc.DeleteComment(r.Context(), uint(id)); err != nil {
		httpx.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
