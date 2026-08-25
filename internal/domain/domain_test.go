package domain

import "testing"

func TestRecordValidationAndTransitions(t *testing.T) {
	record := NewRecord("r1", "祖屋", "家族在此地生活了三代。", "作者", 100, "2025-01-01T00:00:00Z")
	if err := record.Validate(); err != nil {
		t.Fatal(err)
	}
	if err := record.AssignReviewer("审核员", "2025-01-01T01:00:00Z"); err != nil {
		t.Fatal(err)
	}
	if record.Status != StatusReview {
		t.Fatalf("status=%s", record.Status)
	}
	if err := record.UpdateAmount(200, "2025-01-01T02:00:00Z"); err != nil {
		t.Fatal(err)
	}
	if record.Amount != 200 || record.Version != 3 {
		t.Fatalf("record=%+v", record)
	}
	if err := record.AddTag("家乡", "2025-01-01T03:00:00Z"); err != nil {
		t.Fatal(err)
	}
	if !record.RemoveTag("家乡", "2025-01-01T04:00:00Z") {
		t.Fatal("tag was not removed")
	}
}

func TestWorkflowRules(t *testing.T) {
	workflow := NewWorkflow("w1", "r1", "owner", "2025-01-01T00:00:00Z")
	if err := workflow.Advance(StageReview, "ready", "2025-01-01T01:00:00Z"); err != nil {
		t.Fatal(err)
	}
	if err := workflow.Advance(StageConfirm, "approved", "2025-01-01T02:00:00Z"); err != nil {
		t.Fatal(err)
	}
	if err := workflow.Advance(StageArchive, "filed", "2025-01-01T03:00:00Z"); err != nil {
		t.Fatal(err)
	}
	if !workflow.IsComplete() || workflow.LatestNote() != "filed" {
		t.Fatalf("workflow=%+v", workflow)
	}
}

func TestIndexAndRules(t *testing.T) {
	record := NewRecord("r1", "索引故事", "这是一个用于索引的家族故事。", "作者", 1200, "2025-01-01")
	record.Tags = []string{"家书", "家书"}
	index := BuildIndex([]Record{record})
	if index.Size() != 1 || len(index.IDsForTag("家书")) != 1 {
		t.Fatalf("index=%+v", index)
	}
	if ClassifyAmount(record.Amount) != AmountMedium || !CanArchive(Record{Status: StatusConfirmed, Reviewer: "审核"}) {
		t.Fatal("rule result unexpected")
	}
}
