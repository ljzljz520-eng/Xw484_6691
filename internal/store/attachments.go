package store

import (
	"fmt"

	"genealogy-story-organizer/internal/domain"
	"go.etcd.io/bbolt"
)

func (s *Store) SaveAttachment(attachment domain.Attachment) error {
	if err := attachment.Validate(); err != nil {
		return err
	}
	data, err := encode(attachment)
	if err != nil {
		return err
	}
	return s.transactionUpdate(func(tx *bbolt.Tx) error { return tx.Bucket(attachmentBucket).Put(keyFor(attachment.ID), data) })
}

func (s *Store) GetAttachment(id string) (domain.Attachment, error) {
	var attachment domain.Attachment
	err := s.transactionView(func(tx *bbolt.Tx) error {
		value := tx.Bucket(attachmentBucket).Get(keyFor(id))
		if value == nil {
			return fmt.Errorf("attachment %q not found", id)
		}
		return decode(value, &attachment)
	})
	return attachment, err
}

func (s *Store) ListAttachments(recordID string) ([]domain.Attachment, error) {
	result := []domain.Attachment{}
	err := s.transactionView(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(attachmentBucket)
		for _, key := range sortedKeys(bucket) {
			var attachment domain.Attachment
			if err := decode(bucket.Get(key), &attachment); err != nil {
				return err
			}
			if recordID == "" || attachment.RecordID == recordID {
				result = append(result, attachment)
			}
		}
		return nil
	})
	return result, err
}
