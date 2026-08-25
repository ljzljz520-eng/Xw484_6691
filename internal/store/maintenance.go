package store

import (
	"errors"
	"fmt"
	"strings"

	"genealogy-story-organizer/internal/domain"
	"go.etcd.io/bbolt"
)

type IntegrityReport struct {
	Records     int
	Audits      int
	Workflows   int
	Attachments int
	Problems    []string
}

func (s *Store) Inspect() (IntegrityReport, error) {
	report := IntegrityReport{Problems: []string{}}
	if err := s.transactionView(func(tx *bbolt.Tx) error {
		report.Records = tx.Bucket(recordBucket).Stats().KeyN
		report.Audits = tx.Bucket(auditBucket).Stats().KeyN
		report.Workflows = tx.Bucket(workflowBucket).Stats().KeyN
		report.Attachments = tx.Bucket(attachmentBucket).Stats().KeyN
		return nil
	}); err != nil {
		return report, err
	}
	records, err := s.ListRecords()
	if err != nil {
		return report, err
	}
	for _, record := range records {
		if err := record.Validate(); err != nil {
			report.Problems = append(report.Problems, record.ID+": "+err.Error())
		}
	}
	return report, nil
}

func (s *Store) ValidateReferences() ([]string, error) {
	records, err := s.ListRecords()
	if err != nil {
		return nil, err
	}
	known := map[string]bool{}
	for _, record := range records {
		known[record.ID] = true
	}
	problems := []string{}
	workflows, err := s.ListWorkflows()
	if err != nil {
		return nil, err
	}
	for _, workflow := range workflows {
		if !known[workflow.RecordID] {
			problems = append(problems, "workflow="+workflow.ID)
		}
	}
	attachments, err := s.ListAttachments("")
	if err != nil {
		return nil, err
	}
	for _, attachment := range attachments {
		if !known[attachment.RecordID] {
			problems = append(problems, "attachment="+attachment.ID)
		}
	}
	return problems, nil
}

func (s *Store) SaveBundle(record domain.Record, workflow domain.Workflow, attachments []domain.Attachment, audits []domain.AuditEvent) error {
	if err := record.Validate(); err != nil {
		return err
	}
	if workflow.RecordID != record.ID {
		return errors.New("workflow record mismatch")
	}
	if err := workflow.Validate(); err != nil {
		return err
	}
	for _, attachment := range attachments {
		if err := attachment.Validate(); err != nil {
			return err
		}
		if attachment.RecordID != record.ID {
			return errors.New("attachment record mismatch")
		}
	}
	for _, event := range audits {
		if err := event.Validate(); err != nil {
			return err
		}
		if event.RecordID != record.ID {
			return errors.New("audit record mismatch")
		}
	}
	return s.transactionUpdate(func(tx *bbolt.Tx) error {
		recordData, err := encode(record)
		if err != nil {
			return err
		}
		if err := tx.Bucket(recordBucket).Put(keyFor(record.ID), recordData); err != nil {
			return err
		}
		workflowData, err := encode(workflow)
		if err != nil {
			return err
		}
		if err := tx.Bucket(workflowBucket).Put(keyFor(workflow.ID), workflowData); err != nil {
			return err
		}
		for _, attachment := range attachments {
			data, err := encode(attachment)
			if err != nil {
				return err
			}
			if err := tx.Bucket(attachmentBucket).Put(keyFor(attachment.ID), data); err != nil {
				return err
			}
		}
		for _, event := range audits {
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

func (s *Store) Describe() (string, error) {
	report, err := s.Inspect()
	if err != nil {
		return "", err
	}
	if len(report.Problems) > 0 {
		return "invalid: " + strings.Join(report.Problems, ";"), nil
	}
	return fmt.Sprintf("records=%d audits=%d workflows=%d attachments=%d", report.Records, report.Audits, report.Workflows, report.Attachments), nil
}
