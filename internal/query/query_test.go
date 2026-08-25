package query

import (
	"strings"
	"testing"

	"genealogy-story-organizer/internal/domain"
)

func TestSearchAndReport(t *testing.T) {
	records := []domain.Record{
		domain.NewRecord("2", "乙故事", "乙故事的详细叙述。", "作者", 200, "2025-01-01"),
		domain.NewRecord("1", "甲故事", "甲故事的详细叙述。", "作者", 100, "2025-01-01"),
	}
	records[0].Tags = []string{"家书"}
	records[1].Tags = []string{"迁徙"}
	min := int64(150)
	found := Search(records, Filter{MinAmount: &min, Tag: "家书"})
	if len(found) != 1 || found[0].ID != "2" {
		t.Fatalf("found=%v", found)
	}
	report := BuildReport(records, nil)
	if report.TotalAmount != 300 || report.String() == "" {
		t.Fatalf("report=%v", report)
	}
}

func TestExportPartitionAndTimeline(t *testing.T) {
	record := domain.NewRecord("r1", "故事", "这是一个故事文本。", "作者", 40, "2025-01-01")
	record.Status = domain.StatusArchived
	rows := ToExportRows([]domain.Record{record})
	if len(rows) != 1 || !strings.Contains(EncodeCSV([]domain.Record{record}), "故事") {
		t.Fatalf("rows=%v", rows)
	}
	partition := PartitionRecords([]domain.Record{record})
	if !partition.HasArchived() || partition.AmountTotal() != 40 {
		t.Fatalf("partition=%+v", partition)
	}
}
