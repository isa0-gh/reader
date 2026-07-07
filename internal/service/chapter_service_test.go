package service

import (
	"context"
	"errors"
	"testing"

	"github.com/isa0-gh/reader/internal/model"
)

type fakeChapterRepo struct {
	chapters   map[uint]*model.Chapter
	pages      map[uint]*model.S3Object
	nextChID   uint
	nextPageID uint

	lastAddPages []model.S3Object
}

func newFakeChapterRepo() *fakeChapterRepo {
	return &fakeChapterRepo{
		chapters: map[uint]*model.Chapter{},
		pages:    map[uint]*model.S3Object{},
	}
}

func (f *fakeChapterRepo) GetByID(_ context.Context, id uint) (*model.Chapter, error) {
	c, ok := f.chapters[id]
	if !ok {
		return nil, errors.New("chapter not found")
	}
	cp := *c
	for _, p := range f.pages {
		if p.ChapterID != nil && *p.ChapterID == id {
			cp.Pages = append(cp.Pages, *p)
		}
	}
	return &cp, nil
}

func (f *fakeChapterRepo) Create(_ context.Context, c *model.Chapter) error {
	f.nextChID++
	c.ID = f.nextChID
	cp := *c
	f.chapters[c.ID] = &cp
	return nil
}

func (f *fakeChapterRepo) Update(_ context.Context, c *model.Chapter) error {
	existing, ok := f.chapters[c.ID]
	if !ok {
		return errors.New("chapter not found")
	}
	existing.Number = c.Number
	existing.Title = c.Title
	return nil
}

func (f *fakeChapterRepo) AddPages(_ context.Context, pages []model.S3Object) error {
	f.lastAddPages = pages
	for i := range pages {
		f.nextPageID++
		pages[i].ID = f.nextPageID
		cp := pages[i]
		f.pages[cp.ID] = &cp
	}
	return nil
}

func (f *fakeChapterRepo) UpdatePageNumbers(_ context.Context, chapterID uint, numbers map[uint]int) error {
	for pageID, num := range numbers {
		p, ok := f.pages[pageID]
		if !ok || p.ChapterID == nil || *p.ChapterID != chapterID {
			continue
		}
		p.PageNumber = num
	}
	return nil
}

func (f *fakeChapterRepo) DeletePage(_ context.Context, chapterID, pageID uint) error {
	p, ok := f.pages[pageID]
	if ok && p.ChapterID != nil && *p.ChapterID == chapterID {
		delete(f.pages, pageID)
	}
	return nil
}

func (f *fakeChapterRepo) Delete(_ context.Context, id uint) error {
	if _, ok := f.chapters[id]; !ok {
		return errors.New("chapter not found")
	}
	delete(f.chapters, id)
	return nil
}

func TestCreateChapter(t *testing.T) {
	repo := newFakeChapterRepo()
	svc := NewChapterService(repo)

	uploaderID := uint(7)
	created, err := svc.CreateChapter(context.Background(), &model.Chapter{
		SeriesID: 1, UploaderID: &uploaderID, Number: 1, Title: "Ch. 1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created.ID == 0 {
		t.Error("expected an assigned ID")
	}
	if created.UploaderID == nil || *created.UploaderID != uploaderID {
		t.Errorf("uploader id not preserved, got %v", created.UploaderID)
	}
}

func TestUpdateChapter_ReturnsFreshCopy(t *testing.T) {
	repo := newFakeChapterRepo()
	svc := NewChapterService(repo)

	created, _ := svc.CreateChapter(context.Background(), &model.Chapter{SeriesID: 1, Number: 1, Title: "Old"})

	updated, err := svc.UpdateChapter(context.Background(), &model.Chapter{ID: created.ID, Number: 2, Title: "New"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Title != "New" || updated.Number != 2 {
		t.Errorf("got %+v, want Title=New Number=2", updated)
	}
}

func TestUpdateChapter_NotFound(t *testing.T) {
	repo := newFakeChapterRepo()
	svc := NewChapterService(repo)

	if _, err := svc.UpdateChapter(context.Background(), &model.Chapter{ID: 999}); err == nil {
		t.Fatal("expected error for missing chapter, got nil")
	}
}

func TestAddPages_SetsChapterIDOnEveryPage(t *testing.T) {
	repo := newFakeChapterRepo()
	svc := NewChapterService(repo)

	created, _ := svc.CreateChapter(context.Background(), &model.Chapter{SeriesID: 1, Number: 1})

	pages := []model.S3Object{{Key: "a.jpg", PageNumber: 1}, {Key: "b.jpg", PageNumber: 2}}
	if err := svc.AddPages(context.Background(), created.ID, pages); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.lastAddPages) != 2 {
		t.Fatalf("got %d pages passed to repo, want 2", len(repo.lastAddPages))
	}
	for _, p := range repo.lastAddPages {
		if p.ChapterID == nil || *p.ChapterID != created.ID {
			t.Errorf("page %q missing chapter id, got %v want %d", p.Key, p.ChapterID, created.ID)
		}
	}
}

func TestUpdatePageNumbers_Passthrough(t *testing.T) {
	repo := newFakeChapterRepo()
	svc := NewChapterService(repo)

	created, _ := svc.CreateChapter(context.Background(), &model.Chapter{SeriesID: 1, Number: 1})
	svc.AddPages(context.Background(), created.ID, []model.S3Object{{Key: "a.jpg", PageNumber: 1}})

	var pageID uint
	for id := range repo.pages {
		pageID = id
	}

	if err := svc.UpdatePageNumbers(context.Background(), created.ID, map[uint]int{pageID: 5}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.pages[pageID].PageNumber != 5 {
		t.Errorf("got page_number %d, want 5", repo.pages[pageID].PageNumber)
	}
}

func TestDeleteChapter_Passthrough(t *testing.T) {
	repo := newFakeChapterRepo()
	svc := NewChapterService(repo)

	created, _ := svc.CreateChapter(context.Background(), &model.Chapter{SeriesID: 1, Number: 1})
	if err := svc.DeleteChapter(context.Background(), created.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := repo.GetByID(context.Background(), created.ID); err == nil {
		t.Error("expected chapter to be gone after delete")
	}
}
