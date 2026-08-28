package main

import (
	"net/http/httptest"
	"path/filepath"
	"testing"

	"genealogy-story-organizer/internal/application"
	"genealogy-story-organizer/internal/store"
	"genealogy-story-organizer/internal/web"
)

func TestServerEntrypoint(t *testing.T) {
	database, err := store.Open(filepath.Join(t.TempDir(), "entry.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	service, err := application.NewService(database, application.NewSequenceID(), store.StaticClock{Value: "2025-01-01T00:00:00Z"})
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest("GET", "/health", nil)
	response := httptest.NewRecorder()
	web.NewServer(service).Handler().ServeHTTP(response, request)
	if response.Code != 200 {
		t.Fatalf("code=%d", response.Code)
	}
}
