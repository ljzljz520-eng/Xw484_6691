package query

import (
	"sort"
	"strings"

	"genealogy-story-organizer/internal/domain"
)

type TimelineItem struct {
	When   string `json:"when"`
	Action string `json:"action"`
	Actor  string `json:"actor"`
	Note   string `json:"note"`
}

func BuildTimeline(events []domain.AuditEvent) []TimelineItem {
	ordered := domain.SortAudit(events)
	items := make([]TimelineItem, 0, len(ordered))
	for _, event := range ordered {
		items = append(items, TimelineItem{When: event.CreatedAt, Action: event.Action, Actor: event.Actor, Note: event.Note})
	}
	return items
}

func FilterTimeline(items []TimelineItem, actor string) []TimelineItem {
	actor = strings.TrimSpace(actor)
	if actor == "" {
		return append([]TimelineItem(nil), items...)
	}
	result := make([]TimelineItem, 0, len(items))
	for _, item := range items {
		if item.Actor == actor {
			result = append(result, item)
		}
	}
	return result
}

func LatestTimeline(items []TimelineItem) (TimelineItem, bool) {
	if len(items) == 0 {
		return TimelineItem{}, false
	}
	ordered := append([]TimelineItem(nil), items...)
	sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].When < ordered[j].When })
	return ordered[len(ordered)-1], true
}

func CountActions(events []domain.AuditEvent) map[string]int {
	counts := map[string]int{}
	for _, event := range events {
		counts[event.Action]++
	}
	return counts
}
