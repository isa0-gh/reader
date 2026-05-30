package handler

import (
	"encoding/json"
	"net/http"

	"github.com/isa0-gh/reader/internal/config"
)

type ConfigHandler struct{ cfg *config.Config }

func NewConfigHandler(cfg *config.Config) *ConfigHandler { return &ConfigHandler{cfg: cfg} }

func (h *ConfigHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"cdn_url":           h.cfg.CDN,
		"register_disabled": h.cfg.RegisterDisabled,
		"login_disabled":    h.cfg.LoginDisabled,
		"maintenance":       h.cfg.Maintenance,
	})
}
