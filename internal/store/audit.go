package store

import (
	"fmt"

	"genealogy-story-organizer/internal/domain"
	"go.etcd.io/bbolt"
)

func (s *Store) SaveAudit(event domain.AuditEvent) error {
	if err := event.Validate(); err != nil {
		return err
	}
	data, err := encode(event)
	if err != nil {
		return err
	}
	return s.transactionUpdate(func(tx *bbolt.Tx) error { return tx.Bucket(auditBucket).Put(keyFor(event.ID), data) })
}

func (s *Store) ListAudit(recordID string) ([]domain.AuditEvent, error) {
	result := []domain.AuditEvent{}
	err := s.transactionView(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(auditBucket)
		for _, key := range sortedKeys(bucket) {
			var event domain.AuditEvent
			if err := decode(bucket.Get(key), &event); err != nil {
				return err
			}
			if recordID == "" || event.RecordID == recordID {
				result = append(result, event)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return domain.SortAudit(result), nil
}

func (s *Store) NextAuditSequence(recordID string) (int, error) {
	events, err := s.ListAudit(recordID)
	if err != nil {
		return 0, err
	}
	sequence := len(events) + 1
	for _, event := range events {
		if event.Sequence >= sequence {
			sequence = event.Sequence + 1
		}
	}
	return sequence, nil
}

func (s *Store) AuditSummary(recordID string) (string, error) {
	events, err := s.ListAudit(recordID)
	if err != nil {
		return "", err
	}
	if len(events) == 0 {
		return "no audit events", nil
	}
	last := events[len(events)-1]
	return fmt.Sprintf("%d events; last %s", len(events), last.Label()), nil
}
