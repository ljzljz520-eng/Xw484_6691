package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"go.etcd.io/bbolt"
)

var (
	recordBucket     = []byte("records")
	auditBucket      = []byte("audit_events")
	workflowBucket   = []byte("workflows")
	attachmentBucket = []byte("attachments")
)

type Clock interface{ Now() string }

type StaticClock struct{ Value string }

func (c StaticClock) Now() string {
	if c.Value == "" {
		return "2000-01-01T00:00:00Z"
	}
	return c.Value
}

type Store struct {
	db    *bbolt.DB
	mu    sync.RWMutex
	clock Clock
	path  string
}

func Open(path string) (*Store, error) {
	return OpenWithClock(path, StaticClock{})
}

func OpenWithClock(path string, clock Clock) (*Store, error) {
	if path == "" {
		return nil, errors.New("store path is required")
	}
	if clock == nil {
		clock = StaticClock{}
	}
	db, err := bbolt.Open(filepath.Clean(path), 0o600, &bbolt.Options{Timeout: time.Second})
	if err != nil {
		return nil, fmt.Errorf("open bbolt: %w", err)
	}
	store := &Store{db: db, clock: clock, path: filepath.Clean(path)}
	if err := store.initialize(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (s *Store) initialize() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, name := range [][]byte{recordBucket, auditBucket, workflowBucket, attachmentBucket} {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}

func (s *Store) Reopen() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db != nil {
		return errors.New("store is already open")
	}
	db, err := bbolt.Open(s.path, 0o600, &bbolt.Options{Timeout: time.Second})
	if err != nil {
		return err
	}
	s.db = db
	if err := s.initialize(); err != nil {
		_ = db.Close()
		s.db = nil
		return err
	}
	return nil
}

func (s *Store) ensureOpen() error {
	if s.db == nil {
		return errors.New("store is closed")
	}
	return nil
}

func encode(value any) ([]byte, error) { return json.Marshal(value) }

func decode(data []byte, target any) error {
	if len(data) == 0 {
		return errors.New("empty stored value")
	}
	return json.Unmarshal(data, target)
}

func (s *Store) transactionUpdate(fn func(*bbolt.Tx) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.ensureOpen(); err != nil {
		return err
	}
	return s.db.Update(fn)
}

func (s *Store) transactionView(fn func(*bbolt.Tx) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.ensureOpen(); err != nil {
		return err
	}
	return s.db.View(fn)
}

func keyFor(id string) []byte { return []byte(id) }

func sortedKeys(bucket *bbolt.Bucket) [][]byte {
	keys := make([][]byte, 0)
	_ = bucket.ForEach(func(k, v []byte) error {
		if v != nil {
			keys = append(keys, append([]byte(nil), k...))
		}
		return nil
	})
	sort.Slice(keys, func(i, j int) bool { return string(keys[i]) < string(keys[j]) })
	return keys
}

func (s *Store) Path() string { return s.path }

func (s *Store) Now() string { return s.clock.Now() }

func (s *Store) Count(bucketName []byte) (int, error) {
	count := 0
	err := s.transactionView(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return errors.New("bucket not found")
		}
		return bucket.ForEach(func(k, v []byte) error {
			if v != nil {
				count++
			}
			return nil
		})
	})
	return count, err
}
