package store

import (
	"encoding/json"
	"errors"
	"fmt"

	"genealogy-story-organizer/internal/domain"
	"go.etcd.io/bbolt"
)

type Snapshot struct {
	Records     []domain.Record     `json:"records"`
	Audits      []domain.AuditEvent `json:"audits"`
	Workflows   []domain.Workflow   `json:"workflows"`
	Attachments []domain.Attachment `json:"attachments"`
}

func (s *Store) ExportSnapshot() (Snapshot, error) {
	records, err := s.ListRecords()
	if err != nil {
		return Snapshot{}, err
	}
	workflows, err := s.ListWorkflows()
	if err != nil {
		return Snapshot{}, err
	}
	attachments, err := s.ListAttachments("")
	if err != nil {
		return Snapshot{}, err
	}
	audits, err := s.ListAudit("")
	if err != nil {
		return Snapshot{}, err
	}
	return Snapshot{Records: records, Audits: audits, Workflows: workflows, Attachments: attachments}, nil
}

func (s *Store) ExportJSON() ([]byte, error) {
	snapshot, err := s.ExportSnapshot()
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(snapshot, "", "  ")
}

func (s *Store) ImportSnapshot(snapshot Snapshot) error {
	if len(snapshot.Records) == 0 {
		return errors.New("snapshot has no records")
	}
	known := map[string]bool{}
	for _, record := range snapshot.Records {
		if err := record.Validate(); err != nil {
			return err
		}
		if known[record.ID] {
			return fmt.Errorf("duplicate record %s", record.ID)
		}
		known[record.ID] = true
	}
	for _, workflow := range snapshot.Workflows {
		if err := workflow.Validate(); err != nil {
			return err
		}
		if !known[workflow.RecordID] {
			return fmt.Errorf("workflow references unknown record %s", workflow.RecordID)
		}
	}
	return s.transactionUpdate(func(tx *bbolt.Tx) error {
		for _, record := range snapshot.Records {
			data, err := encode(record)
			if err != nil {
				return err
			}
			if err := tx.Bucket(recordBucket).Put(keyFor(record.ID), data); err != nil {
				return err
			}
		}
		for _, workflow := range snapshot.Workflows {
			data, err := encode(workflow)
			if err != nil {
				return err
			}
			if err := tx.Bucket(workflowBucket).Put(keyFor(workflow.ID), data); err != nil {
				return err
			}
		}
		for _, attachment := range snapshot.Attachments {
			if err := attachment.Validate(); err != nil {
				return err
			}
			data, err := encode(attachment)
			if err != nil {
				return err
			}
			if err := tx.Bucket(attachmentBucket).Put(keyFor(attachment.ID), data); err != nil {
				return err
			}
		}
		for _, event := range snapshot.Audits {
			if err := event.Validate(); err != nil {
				return err
			}
			data, err := encode(event)
			if err != nil {
				return err
			}
			if err := tx.Bucket(auditBucket).Put(keyFor(event.ID), data); err != nil {
				return err
			}
		}
		return nil
	})
}
