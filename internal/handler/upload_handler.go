package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/isa0-gh/reader/internal/httpx"
	"github.com/isa0-gh/reader/internal/storage"
)

type UploadHandler struct{ s3 *storage.S3Client }

func NewUploadHandler(s3 *storage.S3Client) *UploadHandler { return &UploadHandler{s3: s3} }

// POST /api/v1/upload/presign
// Body: { "filename": "cover.jpg", "content_type": "image/jpeg", "prefix": "covers" }
// Returns: { "upload_url": "...", "key": "...", "public_url": "..." }
func (h *UploadHandler) Presign(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Filename string `json:"filename"`
		Prefix   string `json:"prefix"` // e.g. "covers", "chapters/1", "avatars"
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Filename == "" {
		httpx.WriteError(w, "filename required", http.StatusBadRequest)
		return
	}

	ext := filepath.Ext(req.Filename)
	key := fmt.Sprintf("%s/%d%s", strings.Trim(req.Prefix, "/"), time.Now().UnixNano(), ext)

	url, err := h.s3.PresignPut(r.Context(), key)
	if err != nil {
		httpx.WriteError(w, "could not generate upload URL", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"upload_url": url,
		"key":        key,
		"public_url": h.s3.PublicURL(key),
		"bucket":     h.s3.Bucket(),
	})
}
