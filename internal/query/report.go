package query

import (
	"fmt"
	"sort"
	"strings"

	"genealogy-story-organizer/internal/domain"
)

type Report struct {
	RecordCount     int              `json:"record_count"`
	AttachmentCount int              `json:"attachment_count"`
	TotalAmount     int64            `json:"total_amount"`
	ByStatus        map[string]int   `json:"by_status"`
	ByCurrency      map[string]int64 `json:"by_currency"`
	Titles          []string         `json:"titles"`
}

func BuildReport(records []domain.Record, attachments []domain.Attachment) Report {
	report := Report{RecordCount: len(records), AttachmentCount: len(attachments), ByStatus: map[string]int{}, ByCurrency: map[string]int64{}, Titles: []string{}}
	for _, record := range records {
		report.TotalAmount += record.Amount
		report.ByStatus[string(record.Status)]++
		report.ByCurrency[record.Currency] += record.Amount
		report.Titles = append(report.Titles, record.Title)
	}
	sort.Strings(report.Titles)
	return report
}

func (r Report) String() string {
	statuses := make([]string, 0, len(r.ByStatus))
	for status, count := range r.ByStatus {
		statuses = append(statuses, fmt.Sprintf("%s=%d", status, count))
	}
	sort.Strings(statuses)
	return fmt.Sprintf("records=%d attachments=%d total=%d statuses=%s", r.RecordCount, r.AttachmentCount, r.TotalAmount, strings.Join(statuses, ","))
}

func (r Report) Empty() bool { return r.RecordCount == 0 && r.AttachmentCount == 0 }

func (r Report) CurrencyTotal(currency string) int64 { return r.ByCurrency[currency] }
