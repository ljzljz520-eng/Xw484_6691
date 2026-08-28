package application

import (
	"path/filepath"
	"testing"

	"genealogy-story-organizer/internal/domain"
	"genealogy-story-organizer/internal/store"
)

func newTestService(t *testing.T) (*Service, *store.Store) {
	t.Helper()
	database, err := store.OpenWithClock(filepath.Join(t.TempDir(), "data.db"), store.StaticClock{Value: "2025-03-01T00:00:00Z"})
	if err != nil {
		t.Fatal(err)
	}
	service, err := NewService(database, NewSequenceID(), store.StaticClock{Value: "2025-03-01T00:00:00Z"})
	if err != nil {
		t.Fatal(err)
	}
	return service, database
}

func TestWorkflowCreateReviewArchive(t *testing.T) {
	service, database := newTestService(t)
	defer database.Close()
	record, err := service.CreateStory("南迁家书", "先祖迁徙后写下了这封家书。", "研究员", 300, []string{"迁徙", "家书"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = service.SubmitForReview(record.ID, "审核员"); err != nil {
		t.Fatal(err)
	}
	if _, err = service.ConfirmReview(record.ID, "审核员", "来源清晰"); err != nil {
		t.Fatal(err)
	}
	archived, err := service.ArchiveStory(record.ID, "档案员")
	if err != nil {
		t.Fatal(err)
	}
	if archived.Status != domain.StatusArchived {
		t.Fatalf("status=%s", archived.Status)
	}
	workflow, err := database.GetWorkflowForRecord(record.ID)
	if err != nil || workflow.Stage != domain.StageArchive {
		t.Fatalf("workflow=%+v err=%v", workflow, err)
	}
	audit, err := database.ListAudit(record.ID)
	if err != nil || len(audit) < 4 {
		t.Fatalf("audit=%v err=%v", audit, err)
	}
}

func TestWorkflowSearchUpdatePublish(t *testing.T) {
	service, database := newTestService(t)
	defer database.Close()
	record, err := service.CreateStory("一门三代", "三代人共同守护族谱和祖屋。", "研究员", 100, nil)
	if err != nil {
		t.Fatal(err)
	}
	updated, err := service.UpdateAmount(record.ID, 175, "编辑")
	if err != nil {
		t.Fatal(err)
	}
	if updated.Amount != 175 {
		t.Fatalf("amount=%d", updated.Amount)
	}
	loaded, err := service.GetStory(record.ID)
	if err != nil || loaded.Amount != 175 {
		t.Fatalf("loaded=%+v err=%v", loaded, err)
	}
}

func TestWorkflowImportReport(t *testing.T) {
	service, database := newTestService(t)
	defer database.Close()
	payload := map[string]any{}
	_ = payload
	result, err := service.ImportPayload(importerPayload())
	if err != nil {
		t.Fatal(err)
	}
	if result.Report.RecordCount != 2 || result.Report.AttachmentCount != 2 {
		t.Fatalf("report=%+v", result.Report)
	}
	if result.Report.TotalAmount != 330 {
		t.Fatalf("total=%d", result.Report.TotalAmount)
	}
}

func TestAuditDigestAndQueue(t *testing.T) {
	service, database := newTestService(t)
	defer database.Close()
	record, err := service.CreateStory("待审核", "这是一条待审核的故事。", "研究员", 50, []string{"待办"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = service.SubmitForReview(record.ID, "审核员"); err != nil {
		t.Fatal(err)
	}
	queue, err := service.ReviewQueue()
	if err != nil || len(queue) != 1 {
		t.Fatalf("queue=%v err=%v", queue, err)
	}
	digest, err := service.AuditDigest(record.ID)
	if err != nil || digest.Events < 2 {
		t.Fatalf("digest=%+v err=%v", digest, err)
	}
}
