package handler

import (
	"encoding/json"
	"net/http"

	"github.com/isa0-gh/reader/internal/config"
)

type ConfigHandler struct{ cfg *config.Config }

func NewConfigHandler(cfg *config.Config) *ConfigHandler { return &ConfigHandler{cfg: cfg} }

type PublicConfigResponse struct {
	CDNURL           string `json:"cdn_url"`
	RegisterDisabled bool   `json:"register_disabled"`
	LoginDisabled    bool   `json:"login_disabled"`
	Maintenance      bool   `json:"maintenance"`
}

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

func (h *ConfigHandler) Get(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, PublicConfigResponse{
		CDNURL:           h.cfg.CDN,
		RegisterDisabled: h.cfg.RegisterDisabled,
		LoginDisabled:    h.cfg.LoginDisabled,
		Maintenance:      h.cfg.Maintenance,
	})
}

func (h *ConfigHandler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, HealthResponse{
		Status:  "ok",
		Service: "reader-api",
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
