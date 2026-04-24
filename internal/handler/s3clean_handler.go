package handler

import (
	"encoding/json"
	"net/http"

	"github.com/isa0-gh/reader/internal/model"
	"github.com/isa0-gh/reader/internal/storage"
	"gorm.io/gorm"
)

type S3CleanHandler struct {
	db *gorm.DB
	s3 *storage.S3Client
}

func NewS3CleanHandler(db *gorm.DB, s3 *storage.S3Client) *S3CleanHandler {
	return &S3CleanHandler{db: db, s3: s3}
}

// orphaned returns S3Objects whose chapter or parent series is soft-deleted or missing.
func (h *S3CleanHandler) orphaned() ([]model.S3Object, error) {
	var objs []model.S3Object
	err := h.db.Unscoped().
		Joins("LEFT JOIN chapters ON chapters.id = s3_objects.chapter_id").
		Joins("LEFT JOIN series ON series.id = chapters.series_id").
		Where("chapters.deleted_at IS NOT NULL OR chapters.id IS NULL OR series.deleted_at IS NOT NULL OR series.id IS NULL").
		Find(&objs).Error
	return objs, err
}

// GET /api/v1/admin/s3/orphaned
func (h *S3CleanHandler) List(w http.ResponseWriter, r *http.Request) {
	objs, err := h.orphaned()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(objs)
}

// DELETE /api/v1/admin/s3/orphaned
func (h *S3CleanHandler) Purge(w http.ResponseWriter, r *http.Request) {
	objs, err := h.orphaned()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var deleted, failed []string
	for _, obj := range objs {
		if err := h.s3.DeleteObject(r.Context(), obj.Key); err != nil {
			failed = append(failed, obj.Key)
			continue
		}
		h.db.Delete(&obj)
		deleted = append(deleted, obj.Key)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"deleted": deleted,
		"failed":  failed,
	})
}
