package store

import (
	"fmt"

	"genealogy-story-organizer/internal/domain"
	"go.etcd.io/bbolt"
)

func (s *Store) SaveWorkflow(workflow domain.Workflow) error {
	if err := workflow.Validate(); err != nil {
		return err
	}
	data, err := encode(workflow)
	if err != nil {
		return err
	}
	return s.transactionUpdate(func(tx *bbolt.Tx) error { return tx.Bucket(workflowBucket).Put(keyFor(workflow.ID), data) })
}

func (s *Store) GetWorkflow(id string) (domain.Workflow, error) {
	var workflow domain.Workflow
	err := s.transactionView(func(tx *bbolt.Tx) error {
		value := tx.Bucket(workflowBucket).Get(keyFor(id))
		if value == nil {
			return fmt.Errorf("workflow %q not found", id)
		}
		return decode(value, &workflow)
	})
	return workflow, err
}

func (s *Store) GetWorkflowForRecord(recordID string) (domain.Workflow, error) {
	workflows, err := s.ListWorkflows()
	if err != nil {
		return domain.Workflow{}, err
	}
	for _, workflow := range workflows {
		if workflow.RecordID == recordID {
			return workflow, nil
		}
	}
	return domain.Workflow{}, fmt.Errorf("workflow for record %q not found", recordID)
}

func (s *Store) ListWorkflows() ([]domain.Workflow, error) {
	result := []domain.Workflow{}
	err := s.transactionView(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(workflowBucket)
		for _, key := range sortedKeys(bucket) {
			var workflow domain.Workflow
			if err := decode(bucket.Get(key), &workflow); err != nil {
				return err
			}
			result = append(result, workflow)
		}
		return nil
	})
	return result, err
}

func (s *Store) AdvanceWorkflow(id string, next domain.WorkflowStage, note string) (domain.Workflow, error) {
	var updated domain.Workflow
	err := s.transactionUpdate(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(workflowBucket)
		value := bucket.Get(keyFor(id))
		if value == nil {
			return fmt.Errorf("workflow %q not found", id)
		}
		if err := decode(value, &updated); err != nil {
			return err
		}
		if err := updated.Advance(next, note, s.Now()); err != nil {
			return err
		}
		data, err := encode(updated)
		if err != nil {
			return err
		}
		return bucket.Put(keyFor(id), data)
	})
	return updated, err
}
