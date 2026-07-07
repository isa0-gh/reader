package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	appmw "github.com/isa0-gh/reader/internal/middleware"
	"github.com/isa0-gh/reader/internal/model"
)

type fakeChapterService struct {
	chapters map[uint]*model.Chapter
}

func (f *fakeChapterService) GetChapter(_ context.Context, id uint) (*model.Chapter, error) {
	c, ok := f.chapters[id]
	if !ok {
		return nil, errors.New("chapter not found")
	}
	return c, nil
}
func (f *fakeChapterService) CreateChapter(_ context.Context, c *model.Chapter) (*model.Chapter, error) {
	return c, nil
}
func (f *fakeChapterService) UpdateChapter(_ context.Context, c *model.Chapter) (*model.Chapter, error) {
	return c, nil
}
func (f *fakeChapterService) AddPages(_ context.Context, _ uint, _ []model.S3Object) error {
	return nil
}
func (f *fakeChapterService) UpdatePageNumbers(_ context.Context, _ uint, _ map[uint]int) error {
	return nil
}
func (f *fakeChapterService) DeletePage(_ context.Context, _, _ uint) error { return nil }
func (f *fakeChapterService) DeleteChapter(_ context.Context, _ uint) error { return nil }

// authorize is the gate added to enforce chapter:update/delete for
// moderators/admins and chapter:update:own/delete:own for the uploader who
// owns the chapter. It's the one place ownership is actually checked, so it
// carries the real security weight of the uploader role.
func TestChapterHandlerAuthorize(t *testing.T) {
	const ownerID, otherID, staffID uint = 1, 2, 99
	chapterID := uint(10)

	svc := &fakeChapterService{chapters: map[uint]*model.Chapter{
		chapterID: {ID: chapterID, UploaderID: ptr(ownerID)},
	}}
	h := NewChapterHandler(svc)

	tests := []struct {
		name       string
		user       *model.User
		action     string
		wantOK     bool
		wantStatus int
	}{
		{"no user in context", nil, "update", false, http.StatusUnauthorized},
		{"reader has no permission", &model.User{ID: staffID, Role: model.RoleReader}, "update", false, http.StatusForbidden},
		{"uploader who does not own the chapter", &model.User{ID: otherID, Role: model.RoleUploader}, "update", false, http.StatusForbidden},
		{"uploader who owns the chapter", &model.User{ID: ownerID, Role: model.RoleUploader}, "update", true, 0},
		{"moderator has blanket update", &model.User{ID: staffID, Role: model.RoleModerator}, "update", true, 0},
		{"admin has blanket delete", &model.User{ID: staffID, Role: model.RoleAdmin}, "delete", true, 0},
		{"owning uploader can delete", &model.User{ID: ownerID, Role: model.RoleUploader}, "delete", true, 0},
		{"non-owning uploader cannot delete", &model.User{ID: otherID, Role: model.RoleUploader}, "delete", false, http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.user != nil {
				ctx = context.WithValue(ctx, appmw.UserContextKey, tt.user)
			}
			req := httptest.NewRequest(http.MethodPut, "/chapters/10", nil).WithContext(ctx)
			w := httptest.NewRecorder()

			_, ok := h.authorize(w, req, chapterID, tt.action)
			if ok != tt.wantOK {
				t.Fatalf("authorize() ok = %v, want %v", ok, tt.wantOK)
			}
			if !ok && w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestChapterHandlerAuthorize_ChapterNotFound(t *testing.T) {
	svc := &fakeChapterService{chapters: map[uint]*model.Chapter{}}
	h := NewChapterHandler(svc)

	user := &model.User{ID: 1, Role: model.RoleAdmin}
	ctx := context.WithValue(context.Background(), appmw.UserContextKey, user)
	req := httptest.NewRequest(http.MethodPut, "/chapters/404", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	if _, ok := h.authorize(w, req, 404, "update"); ok {
		t.Fatal("expected ok=false for a chapter that doesn't exist")
	}
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func ptr(u uint) *uint { return &u }
