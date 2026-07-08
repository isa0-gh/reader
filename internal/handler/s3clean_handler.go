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

// orphaned returns S3Objects whose chapter/parent series is soft-deleted or
// missing, or whose owning series/user/background (for covers/avatars/site
// wallpapers) is soft-deleted, hard-deleted, or missing.
func (h *S3CleanHandler) orphaned() ([]model.S3Object, error) {
	// Initialized (not nil) so json.Marshal always encodes "[]", never
	// "null" — the frontend expects an array to call .length on unconditionally.
	objs := []model.S3Object{}

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

	// 2. Covers/avatars/backgrounds orphaned by series/user deletion, or just
	// unreferenced. All three are S3Objects with chapter_id NULL, so they
	// share this query — each new "S3Object with no chapter_id" use case
	// (avatars, now backgrounds) needs its own join here, or it looks
	// unreferenced by everything else and gets purged. Backgrounds have no
	// soft-delete, so a plain "no matching row" check is enough — no
	// deleted_at column to also check like series/users.
	var covers []model.S3Object
	if err := h.db.Unscoped().
		Table("s3_objects").
		Joins("LEFT JOIN series ON series.cover_image_id = s3_objects.id").
		Joins("LEFT JOIN users ON users.avatar_id = s3_objects.id").
		Joins("LEFT JOIN backgrounds ON backgrounds.image_id = s3_objects.id").
		Where("s3_objects.chapter_id IS NULL").
		Where("series.deleted_at IS NOT NULL OR series.id IS NULL").
		Where("users.deleted_at IS NOT NULL OR users.id IS NULL").
		Where("backgrounds.id IS NULL").
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

	deleted, failed := []string{}, []string{}

	for _, obj := range objs {
		if err := h.s3.DeleteObject(r.Context(), obj.Bucket, obj.Key); err != nil {
			log.Printf("s3clean: failed to delete s3://%s/%s: %v", obj.Bucket, obj.Key, err)
			failed = append(failed, obj.Key)
			continue
		}
		// If this was a cover image or avatar still referenced by a soft-deleted
		// series/user, the fk_series_cover_image / fk_users_avatar constraints
		// (ON DELETE SET NULL) null out that reference automatically instead of
		// blocking the delete.
		if err := h.db.Unscoped().Delete(&obj).Error; err != nil {
			log.Printf("s3clean: deleted s3://%s/%s but failed to remove DB row: %v", obj.Bucket, obj.Key, err)
			failed = append(failed, obj.Key)
			continue
		}
		deleted = append(deleted, obj.Key)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"deleted": deleted,
		"failed":  failed,
	})
}
