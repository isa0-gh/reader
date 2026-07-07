package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/isa0-gh/reader/internal/model"
)

type fakeSeriesRepo struct {
	series    map[uint]*model.Series
	nextID    uint
	createErr error
	updateErr error
	getErr    error
	deleteErr error
}

func newFakeSeriesRepo() *fakeSeriesRepo {
	return &fakeSeriesRepo{series: map[uint]*model.Series{}}
}

func (f *fakeSeriesRepo) GetByID(_ context.Context, id uint) (*model.Series, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	s, ok := f.series[id]
	if !ok {
		return nil, errors.New("series not found")
	}
	cp := *s
	return &cp, nil
}

func (f *fakeSeriesRepo) Create(_ context.Context, s *model.Series) error {
	if f.createErr != nil {
		return f.createErr
	}
	f.nextID++
	s.ID = f.nextID
	cp := *s
	f.series[s.ID] = &cp
	return nil
}

func (f *fakeSeriesRepo) Update(_ context.Context, s *model.Series) error {
	if f.updateErr != nil {
		return f.updateErr
	}
	if _, ok := f.series[s.ID]; !ok {
		return errors.New("series not found")
	}
	cp := *s
	f.series[s.ID] = &cp
	return nil
}

func (f *fakeSeriesRepo) List(_ context.Context, limit, offset int, q string) ([]model.Series, error) {
	list := make([]model.Series, 0, len(f.series))
	for _, s := range f.series {
		if q != "" && !strings.Contains(strings.ToLower(s.Title), strings.ToLower(q)) {
			continue
		}
		list = append(list, *s)
	}
	if offset > len(list) {
		return []model.Series{}, nil
	}
	list = list[offset:]
	if limit > 0 && limit < len(list) {
		list = list[:limit]
	}
	return list, nil
}

func (f *fakeSeriesRepo) Count(_ context.Context, q string) (int64, error) {
	var count int64
	for _, s := range f.series {
		if q != "" && !strings.Contains(strings.ToLower(s.Title), strings.ToLower(q)) {
			continue
		}
		count++
	}
	return count, nil
}

func (f *fakeSeriesRepo) Delete(_ context.Context, id uint) error {
	if f.deleteErr != nil {
		return f.deleteErr
	}
	if _, ok := f.series[id]; !ok {
		return errors.New("series not found")
	}
	delete(f.series, id)
	return nil
}

func TestCreateSeries(t *testing.T) {
	repo := newFakeSeriesRepo()
	svc := NewSeriesService(repo)

	created, err := svc.CreateSeries(context.Background(), &model.Series{Title: "One Piece", Slug: "one-piece"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created.ID == 0 {
		t.Error("expected an assigned ID")
	}
	if created.Title != "One Piece" {
		t.Errorf("got title %q, want %q", created.Title, "One Piece")
	}
}

func TestCreateSeries_PropagatesRepoError(t *testing.T) {
	repo := newFakeSeriesRepo()
	repo.createErr = errors.New("db down")
	svc := NewSeriesService(repo)

	if _, err := svc.CreateSeries(context.Background(), &model.Series{Title: "X", Slug: "x"}); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUpdateSeries_ReturnsFreshCopy(t *testing.T) {
	repo := newFakeSeriesRepo()
	svc := NewSeriesService(repo)

	created, err := svc.CreateSeries(context.Background(), &model.Series{Title: "Old", Slug: "old"})
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	updated, err := svc.UpdateSeries(context.Background(), &model.Series{ID: created.ID, Title: "New", Slug: "new"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Title != "New" {
		t.Errorf("got title %q, want %q", updated.Title, "New")
	}
}

func TestUpdateSeries_NotFound(t *testing.T) {
	repo := newFakeSeriesRepo()
	svc := NewSeriesService(repo)

	if _, err := svc.UpdateSeries(context.Background(), &model.Series{ID: 999, Title: "X"}); err == nil {
		t.Fatal("expected error for missing series, got nil")
	}
}

func TestDeleteSeries_Passthrough(t *testing.T) {
	repo := newFakeSeriesRepo()
	svc := NewSeriesService(repo)

	created, _ := svc.CreateSeries(context.Background(), &model.Series{Title: "X", Slug: "x"})
	if err := svc.DeleteSeries(context.Background(), created.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := repo.GetByID(context.Background(), created.ID); err == nil {
		t.Error("expected series to be gone after delete")
	}
}

func TestListSeries_Passthrough(t *testing.T) {
	repo := newFakeSeriesRepo()
	svc := NewSeriesService(repo)

	svc.CreateSeries(context.Background(), &model.Series{Title: "A", Slug: "a"})
	svc.CreateSeries(context.Background(), &model.Series{Title: "B", Slug: "b"})

	list, total, err := svc.ListSeries(context.Background(), 0, 0, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 2 || total != 2 {
		t.Errorf("got %d series (total=%d), want 2 (total=2)", len(list), total)
	}
}

func TestListSeries_SearchFiltersByTitle(t *testing.T) {
	repo := newFakeSeriesRepo()
	svc := NewSeriesService(repo)

	svc.CreateSeries(context.Background(), &model.Series{Title: "One Piece", Slug: "one-piece"})
	svc.CreateSeries(context.Background(), &model.Series{Title: "Naruto", Slug: "naruto"})

	list, total, err := svc.ListSeries(context.Background(), 0, 0, "one")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 1 || len(list) != 1 || list[0].Title != "One Piece" {
		t.Errorf("got %+v (total=%d), want only One Piece (total=1)", list, total)
	}
}

func TestListSeries_LimitAndOffset(t *testing.T) {
	repo := newFakeSeriesRepo()
	svc := NewSeriesService(repo)

	for i := 0; i < 5; i++ {
		svc.CreateSeries(context.Background(), &model.Series{Title: "S", Slug: "s"})
	}

	list, total, err := svc.ListSeries(context.Background(), 2, 3, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 5 {
		t.Errorf("got total %d, want 5 (total should ignore limit/offset)", total)
	}
	if len(list) != 2 {
		t.Errorf("got %d series, want 2 (limit)", len(list))
	}
}
