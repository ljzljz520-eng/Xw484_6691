package integration

import (
	"path/filepath"
	"testing"

	"genealogy-story-organizer/internal/application"
	"genealogy-story-organizer/internal/store"
)

func TestLifecycleSnapshot(t *testing.T) {
	database, err := store.OpenWithClock(filepath.Join(t.TempDir(), "integration.db"), store.StaticClock{Value: "2025-04-01T00:00:00Z"})
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	service, err := application.NewService(database, application.NewSequenceID(), store.StaticClock{Value: "2025-04-01T00:00:00Z"})
	if err != nil {
		t.Fatal(err)
	}
	record, err := service.CreateStory("归档故事", "这个故事将经过完整流程归档。", "作者", 77, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = service.SubmitForReview(record.ID, "审核"); err != nil {
		t.Fatal(err)
	}
	if _, err = service.ConfirmReview(record.ID, "审核", "确认"); err != nil {
		t.Fatal(err)
	}
	if _, err = service.ArchiveStory(record.ID, "档案"); err != nil {
		t.Fatal(err)
	}
	snapshot, err := BuildSnapshot(service, record.ID)
	if err != nil || !snapshot.Complete() {
		t.Fatalf("snapshot=%+v err=%v", snapshot, err)
	}
}
