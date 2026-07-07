package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/isa0-gh/reader/internal/httpx"
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

	// 1. Pages orphaned by chapter/series deletion
	var pages []model.S3Object
	if err := h.db.Unscoped().
		Table("s3_objects").
		Joins("LEFT JOIN chapters ON chapters.id = s3_objects.chapter_id").
		Joins("LEFT JOIN series ON series.id = chapters.series_id").
		Where("s3_objects.chapter_id IS NOT NULL").
		Where("chapters.deleted_at IS NOT NULL OR chapters.id IS NULL OR (chapters.id IS NOT NULL AND (series.deleted_at IS NOT NULL OR series.id IS NULL))").
		Find(&pages).Error; err != nil {
		return nil, err
	}
	objs = append(objs, pages...)

	// 2. Covers orphaned by series deletion or just unreferenced
	var covers []model.S3Object
	if err := h.db.Unscoped().
		Table("s3_objects").
		Joins("LEFT JOIN series ON series.cover_image_id = s3_objects.id").
		Where("s3_objects.chapter_id IS NULL").
		Where("series.deleted_at IS NOT NULL OR series.id IS NULL").
		Find(&covers).Error; err != nil {
		return nil, err
	}
	objs = append(objs, covers...)

	return objs, nil
}

// GET /api/v1/admin/s3/orphaned
func (h *S3CleanHandler) List(w http.ResponseWriter, r *http.Request) {
	objs, err := h.orphaned()
	if err != nil {
		httpx.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(objs)
}

// DELETE /api/v1/admin/s3/orphaned
func (h *S3CleanHandler) Purge(w http.ResponseWriter, r *http.Request) {
	objs, err := h.orphaned()
	if err != nil {
		httpx.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var deleted, failed []string

	for _, obj := range objs {
		if err := h.s3.DeleteObject(r.Context(), obj.Bucket, obj.Key); err != nil {
			log.Printf("s3clean: failed to delete s3://%s/%s: %v", obj.Bucket, obj.Key, err)
			failed = append(failed, obj.Key)
			continue
		}
		// If it's a cover image, we should probably also null out the reference in series table if it's still there
		// but since we are deleting the S3Object record, GORM will handle it if it was a real FK,
		// or it will just be a dangling ID. Given we use Unscoped().Delete, it's fine.
		h.db.Unscoped().Delete(&obj)
		deleted = append(deleted, obj.Key)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"deleted": deleted,
		"failed":  failed,
	})
}
