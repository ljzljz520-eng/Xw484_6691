package store

import (
	"path/filepath"
	"testing"

	"genealogy-story-organizer/internal/domain"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "genealogy.db")
	database, err := OpenWithClock(path, StaticClock{Value: "2025-02-01T00:00:00Z"})
	if err != nil {
		t.Fatal(err)
	}
	record := domain.NewRecord("r1", "祖谱", "这是一段可重开的家族故事。", "研究员", 88, "2025-02-01T00:00:00Z")
	if err := database.SaveRecord(record); err != nil {
		t.Fatal(err)
	}
	workflow := domain.NewWorkflow("w1", record.ID, "研究员", "2025-02-01T00:00:00Z")
	if err := database.SaveWorkflow(workflow); err != nil {
		t.Fatal(err)
	}
	if err := database.Close(); err != nil {
		t.Fatal(err)
	}
	if err := database.Reopen(); err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	loaded, err := database.GetRecord("r1")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Title != record.Title || loaded.Amount != 88 {
		t.Fatalf("loaded=%+v", loaded)
	}
}

func TestStoreListsEntities(t *testing.T) {
	database, err := OpenWithClock(filepath.Join(t.TempDir(), "data.db"), StaticClock{Value: "2025-02-01T00:00:00Z"})
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	record := domain.NewRecord("r1", "旧宅", "家中旧宅的来历记录。", "研究员", 10, "2025-02-01T00:00:00Z")
	if err := database.SaveRecord(record); err != nil {
		t.Fatal(err)
	}
	if err := database.SaveAttachment(domain.NewAttachment("a1", "r1", "来源.txt", "text/plain", "abc", "2025-02-01T00:00:00Z", 3)); err != nil {
		t.Fatal(err)
	}
	if got, err := database.ListAttachments("r1"); err != nil || len(got) != 1 {
		t.Fatalf("attachments=%v err=%v", got, err)
	}
	if got, err := database.Count(recordBucket); err != nil || got != 1 {
		t.Fatalf("count=%d err=%v", got, err)
	}
}

func TestSnapshotRoundTrip(t *testing.T) {
	database, err := OpenWithClock(filepath.Join(t.TempDir(), "snapshot.db"), StaticClock{Value: "2025-02-01T00:00:00Z"})
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	record := domain.NewRecord("r1", "快照故事", "快照需要保留完整业务记录。", "研究员", 66, "2025-02-01T00:00:00Z")
	workflow := domain.NewWorkflow("w1", "r1", "研究员", "2025-02-01T00:00:00Z")
	if err := database.SaveBundle(record, workflow, nil, nil); err != nil {
		t.Fatal(err)
	}
	snapshot, err := database.ExportSnapshot()
	if err != nil || len(snapshot.Records) != 1 {
		t.Fatalf("snapshot=%+v err=%v", snapshot, err)
	}
	if len(snapshot.Records) == 0 {
		t.Fatal("empty snapshot")
	}
	if err := database.ImportSnapshot(snapshot); err != nil {
		t.Fatal(err)
	}
	if data, err := database.ExportJSON(); err != nil || len(data) == 0 {
		t.Fatalf("json=%q err=%v", data, err)
	}
}
