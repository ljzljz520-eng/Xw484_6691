package application

import (
	"errors"
	"fmt"
	"strings"

	"genealogy-story-organizer/internal/domain"
	"genealogy-story-organizer/internal/store"
)

type IDSource interface{ Next(prefix string) string }

type SequenceID struct{ counters map[string]int }

func NewSequenceID() *SequenceID { return &SequenceID{counters: map[string]int{}} }

func (s *SequenceID) Next(prefix string) string {
	s.counters[prefix]++
	return fmt.Sprintf("%s-%03d", prefix, s.counters[prefix])
}

type Service struct {
	store       *store.Store
	ids         IDSource
	clock       store.Clock
	amountCache int64
}

func NewService(database *store.Store, ids IDSource, clock store.Clock) (*Service, error) {
	if database == nil {
		return nil, errors.New("store is required")
	}
	if ids == nil {
		ids = NewSequenceID()
	}
	if clock == nil {
		clock = store.StaticClock{}
	}
	return &Service{store: database, ids: ids, clock: clock}, nil
}

func (s *Service) CreateStory(title, narrative, author string, amount int64, tags []string) (domain.Record, error) {
	now := s.clock.Now()
	record := domain.NewRecord(s.ids.Next("record"), title, narrative, author, amount, now)
	record.Tags = domain.NormalizeTags(tags)
	if err := record.Validate(); err != nil {
		return domain.Record{}, err
	}
	workflow := domain.NewWorkflow(s.ids.Next("workflow"), record.ID, author, now)
	if err := s.store.SaveRecord(record); err != nil {
		return domain.Record{}, err
	}
	if err := s.store.SaveWorkflow(workflow); err != nil {
		return domain.Record{}, err
	}
	if err := s.writeAudit(record.ID, "created", author, "story registered"); err != nil {
		return domain.Record{}, err
	}
	return record, nil
}

func (s *Service) Attach(recordID, name, mediaType, checksum string, size int64) (domain.Attachment, error) {
	if _, err := s.store.GetRecord(recordID); err != nil {
		return domain.Attachment{}, err
	}
	attachment := domain.NewAttachment(s.ids.Next("attachment"), recordID, name, mediaType, checksum, s.clock.Now(), size)
	if err := s.store.SaveAttachment(attachment); err != nil {
		return domain.Attachment{}, err
	}
	if err := s.writeAudit(recordID, "attachment-added", "system", attachment.DisplayName()); err != nil {
		return domain.Attachment{}, err
	}
	return attachment, nil
}

func (s *Service) writeAudit(recordID, action, actor, note string) error {
	sequence, err := s.store.NextAuditSequence(recordID)
	if err != nil {
		return err
	}
	event := domain.NewAuditEvent(s.ids.Next("audit"), recordID, action, actor, note, s.clock.Now(), sequence)
	return s.store.SaveAudit(event)
}

func (s *Service) auditMutation(recordID, action, actor, note string) error {
	if strings.TrimSpace(actor) == "" {
		actor = "system"
	}
	return s.writeAudit(recordID, action, actor, note)
}

func (s *Service) Database() *store.Store { return s.store }

func (s *Service) ListStories() ([]domain.Record, error) { return s.store.ListRecords() }

func (s *Service) GetStory(id string) (domain.Record, error) { return s.store.GetRecord(id) }
