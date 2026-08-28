package web

import (
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"genealogy-story-organizer/internal/application"
	"genealogy-story-organizer/internal/store"
)

func TestWebAssets(t *testing.T) {
	database, err := store.Open(filepath.Join(t.TempDir(), "web.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	service, err := application.NewService(database, application.NewSequenceID(), store.StaticClock{Value: "2025-01-01T00:00:00Z"})
	if err != nil {
		t.Fatal(err)
	}
	server := NewServer(service)
	request := httptest.NewRequest("GET", "/health", nil)
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != 200 || !strings.Contains(response.Body.String(), "ok") {
		t.Fatalf("code=%d body=%s", response.Code, response.Body.String())
	}
}
