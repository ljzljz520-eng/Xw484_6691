package query

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"genealogy-story-organizer/internal/domain"
)

type ExportRow struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Status   string `json:"status"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	Tags     string `json:"tags"`
}

func ToExportRows(records []domain.Record) []ExportRow {
	ordered := domain.SortRecords(records, false)
	rows := make([]ExportRow, 0, len(ordered))
	for _, record := range ordered {
		tags := append([]string(nil), record.Tags...)
		sort.Strings(tags)
		rows = append(rows, ExportRow{ID: record.ID, Title: record.Title, Status: string(record.Status), Amount: record.Amount, Currency: record.Currency, Tags: strings.Join(tags, ",")})
	}
	return rows
}

func EncodeJSON(records []domain.Record) ([]byte, error) {
	if records == nil {
		return nil, errors.New("records cannot be nil")
	}
	return json.MarshalIndent(ToExportRows(records), "", "  ")
}

func EncodeCSV(records []domain.Record) string {
	rows := ToExportRows(records)
	lines := []string{"id,title,status,amount,currency,tags"}
	for _, row := range rows {
		lines = append(lines, strings.Join([]string{csvField(row.ID), csvField(row.Title), csvField(row.Status), fmt.Sprintf("%d", row.Amount), csvField(row.Currency), csvField(row.Tags)}, ","))
	}
	return strings.Join(lines, "\n")
}

func csvField(value string) string {
	if strings.ContainsAny(value, ",\"\n") {
		return "\"" + strings.ReplaceAll(value, "\"", "\"\"") + "\""
	}
	return value
}

func ParseStatus(value string) (domain.RecordStatus, error) {
	status := domain.RecordStatus(strings.ToLower(strings.TrimSpace(value)))
	if !status.Valid() {
		return "", fmt.Errorf("invalid status %q", value)
	}
	return status, nil
}

func MatchAll(records []domain.Record, filters []Filter) []domain.Record {
	if len(filters) == 0 {
		return Search(records, Filter{})
	}
	current := records
	for _, filter := range filters {
		current = Search(current, filter)
	}
	return current
}
