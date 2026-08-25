package flow023

import (
	"strings"

	"genealogy-story-organizer/internal/domain"
	"genealogy-story-organizer/internal/query"
)

func Summarize(records []domain.Record) string {
	rows := query.ToExportRows(records)
	parts := make([]string, 0, len(rows))
	for _, row := range rows {
		parts = append(parts, row.ID+":"+row.Title+":"+row.Status)
	}
	return strings.Join(parts, "|")
}

func Amounts(records []domain.Record) map[string]domain.AmountBand {
	result := map[string]domain.AmountBand{}
	for _, record := range records {
		result[record.ID] = domain.ClassifyAmount(record.Amount)
	}
	return result
}
