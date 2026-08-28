package store

import (
	"errors"
	"fmt"

	"genealogy-story-organizer/internal/domain"
	"go.etcd.io/bbolt"
)

func (s *Store) SaveRecord(record domain.Record) error {
	if err := record.Validate(); err != nil {
		return err
	}
	data, err := encode(record)
	if err != nil {
		return fmt.Errorf("encode record: %w", err)
	}
	return s.transactionUpdate(func(tx *bbolt.Tx) error {
		return tx.Bucket(recordBucket).Put(keyFor(record.ID), data)
	})
}

func (s *Store) GetRecord(id string) (domain.Record, error) {
	var record domain.Record
	err := s.transactionView(func(tx *bbolt.Tx) error {
		value := tx.Bucket(recordBucket).Get(keyFor(id))
		if value == nil {
			return fmt.Errorf("record %q not found", id)
		}
		return decode(value, &record)
	})
	return record, err
}

func (s *Store) ListRecords() ([]domain.Record, error) {
	result := []domain.Record{}
	err := s.transactionView(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(recordBucket)
		for _, key := range sortedKeys(bucket) {
			var record domain.Record
			if err := decode(bucket.Get(key), &record); err != nil {
				return err
			}
			result = append(result, record)
		}
		return nil
	})
	return result, err
}

func (s *Store) DeleteRecord(id string) error {
	return s.transactionUpdate(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(recordBucket)
		if bucket.Get(keyFor(id)) == nil {
			return errors.New("record not found")
		}
		return bucket.Delete(keyFor(id))
	})
}

func (s *Store) UpdateRecord(id string, update func(*domain.Record) error) (domain.Record, error) {
	var updated domain.Record
	err := s.transactionUpdate(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(recordBucket)
		value := bucket.Get(keyFor(id))
		if value == nil {
			return fmt.Errorf("record %q not found", id)
		}
		if err := decode(value, &updated); err != nil {
			return err
		}
		if err := update(&updated); err != nil {
			return err
		}
		if err := updated.Validate(); err != nil {
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
