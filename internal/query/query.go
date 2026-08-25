package query

import (
	"sort"
	"strings"

	"genealogy-story-organizer/internal/domain"
)

type Filter struct {
	Text            string
	Status          domain.RecordStatus
	MinAmount       *int64
	MaxAmount       *int64
	Tag             string
	IncludeRejected bool
}

func Search(records []domain.Record, filter Filter) []domain.Record {
	result := make([]domain.Record, 0, len(records))
	for _, record := range records {
		if !matches(record, filter) {
			continue
		}
		result = append(result, record.Clone())
	}
	sort.SliceStable(result, func(i, j int) bool {
		if strings.ToLower(result[i].Title) == strings.ToLower(result[j].Title) {
			return result[i].ID < result[j].ID
		}
		return strings.ToLower(result[i].Title) < strings.ToLower(result[j].Title)
	})
	return result
}

func matches(record domain.Record, filter Filter) bool {
	if !filter.IncludeRejected && !record.IsVisible() {
		return false
	}
	if filter.Status != "" && record.Status != filter.Status {
		return false
	}
	if filter.MinAmount != nil && record.Amount < *filter.MinAmount {
		return false
	}
	if filter.MaxAmount != nil && record.Amount > *filter.MaxAmount {
		return false
	}
	if filter.Text != "" {
		text := strings.ToLower(filter.Text)
		if !strings.Contains(strings.ToLower(record.Title), text) && !strings.Contains(strings.ToLower(record.Narrative), text) {
			return false
		}
	}
	if filter.Tag != "" {
		found := false
		for _, tag := range record.Tags {
			if strings.EqualFold(tag, filter.Tag) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func GroupByStatus(records []domain.Record) map[domain.RecordStatus][]domain.Record {
	groups := map[domain.RecordStatus][]domain.Record{}
	for _, record := range records {
		groups[record.Status] = append(groups[record.Status], record.Clone())
	}
	return groups
}

func TotalsByCurrency(records []domain.Record) map[string]int64 {
	totals := map[string]int64{}
	for _, record := range records {
		totals[record.Currency] += record.Amount
	}
	return totals
}
