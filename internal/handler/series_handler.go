package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/isa0-gh/reader/internal/service"
)

type SeriesHandler struct{ svc service.SeriesService }

func NewSeriesHandler(svc service.SeriesService) *SeriesHandler {
	return &SeriesHandler{svc: svc}
}

func (h *SeriesHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	series, err := h.svc.GetSeries(r.Context(), uint(id))
	if err != nil {
		http.Error(w, "series not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(series)
}
