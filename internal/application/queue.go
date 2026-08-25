package application

import (
	"errors"
	"sort"
	"strings"

	"genealogy-story-organizer/internal/domain"
	"genealogy-story-organizer/internal/query"
)

type ReviewQueueItem struct {
	Record domain.Record
	Events []domain.AuditEvent
	Score  int
}

func (s *Service) ReviewQueue() ([]ReviewQueueItem, error) {
	records, err := s.store.ListRecords()
	if err != nil {
		return nil, err
	}
	queue := []ReviewQueueItem{}
	for _, record := range records {
		if record.Status != domain.StatusReview {
			continue
		}
		events, eventErr := s.store.ListAudit(record.ID)
		if eventErr != nil {
			return nil, eventErr
		}
		queue = append(queue, ReviewQueueItem{Record: record, Events: events, Score: reviewScore(record, events)})
	}
	sort.SliceStable(queue, func(i, j int) bool {
		if queue[i].Score == queue[j].Score {
			return queue[i].Record.ID < queue[j].Record.ID
		}
		return queue[i].Score > queue[j].Score
	})
	return queue, nil
}

func reviewScore(record domain.Record, events []domain.AuditEvent) int {
	score := len(query.FilterTimeline(query.BuildTimeline(events), ""))
	if record.Reviewer != "" {
		score += 2
	}
	if len(record.Tags) > 0 {
		score++
	}
	if record.Amount >= 10000 {
		score += 3
	}
	return score
}

func (s *Service) PublishChange(recordID, actor string) (string, error) {
	record, err := s.store.GetRecord(recordID)
	if err != nil {
		return "", err
	}
	if !record.IsVisible() {
		return "", errors.New("rejected records cannot be published")
	}
	actor = domain.NormalizeActor(actor)
	if err := s.auditMutation(recordID, "published", actor, record.Summary()); err != nil {
		return "", err
	}
	return strings.TrimSpace(record.Summary()), nil
}

func (s *Service) Timeline(recordID string) ([]query.TimelineItem, error) {
	events, err := s.store.ListAudit(recordID)
	if err != nil {
		return nil, err
	}
	return query.BuildTimeline(events), nil
}
