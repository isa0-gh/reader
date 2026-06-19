package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/isa0-gh/reader/internal/config"
	"github.com/isa0-gh/reader/internal/handler"
)

func TestConfigHandlerGetReturnsOnlyPublicConfig(t *testing.T) {
	h := handler.NewConfigHandler(&config.Config{
		Port:             "8080",
		DBDSN:            "postgres://secret-dsn",
		Debug:            true,
		CDN:              "https://cdn.example.com",
		RegisterDisabled: true,
		LoginDisabled:    false,
		Maintenance:      true,
		S3: config.S3Config{
			Endpoint:        "https://s3.example.com",
			Region:          "auto",
			AccessKeyID:     "secret-key-id",
			SecretAccessKey: "secret-access-key",
			Bucket:          "reader-storage",
			UsePathStyle:    true,
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/config", nil)
	rec := httptest.NewRecorder()

	h.Get(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if contentType := rec.Header().Get("Content-Type"); !strings.HasPrefix(contentType, "application/json") {
		t.Fatalf("expected JSON content type, got %q", contentType)
	}

	rawBody := rec.Body.String()
	for _, secretValue := range []string{"secret-dsn", "secret-key-id", "secret-access-key", "reader-storage"} {
		if strings.Contains(rawBody, secretValue) {
			t.Fatalf("response exposed secret value %q", secretValue)
		}
	}

	var body map[string]any
	if err := json.NewDecoder(strings.NewReader(rawBody)).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	want := map[string]any{
		"cdn_url":           "https://cdn.example.com",
		"register_disabled": true,
		"login_disabled":    false,
		"maintenance":       true,
	}
	if len(body) != len(want) {
		t.Fatalf("expected only public config keys %v, got %v", want, body)
	}
	for key, value := range want {
		if body[key] != value {
			t.Fatalf("expected %s=%v, got %v", key, value, body[key])
		}
	}

	for _, secretKey := range []string{
		"port",
		"db_dsn",
		"debug",
		"s3",
		"s3_endpoint",
		"s3_access_key_id",
		"s3_secret_access_key",
		"s3_bucket",
	} {
		if _, ok := body[secretKey]; ok {
			t.Fatalf("response exposed non-public config key %q", secretKey)
		}
	}
}

func TestConfigHandlerHealth(t *testing.T) {
	h := handler.NewConfigHandler(&config.Config{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rec := httptest.NewRecorder()

	h.Health(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("expected status ok, got %q", body["status"])
	}
	if body["service"] != "reader-api" {
		t.Fatalf("expected service reader-api, got %q", body["service"])
	}
}
