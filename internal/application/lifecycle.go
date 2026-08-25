package application

import (
	"errors"
	"fmt"

	"genealogy-story-organizer/internal/domain"
)

func (s *Service) SubmitForReview(recordID, reviewer string) (domain.Record, error) {
	record, err := s.store.UpdateRecord(recordID, func(target *domain.Record) error {
		return target.AssignReviewer(reviewer, s.clock.Now())
	})
	if err != nil {
		return domain.Record{}, err
	}
	workflow, err := s.store.GetWorkflowForRecord(recordID)
	if err != nil {
		return domain.Record{}, err
	}
	if _, err := s.store.AdvanceWorkflow(workflow.ID, domain.StageReview, "submitted by "+reviewer); err != nil {
		return domain.Record{}, err
	}
	if err := s.auditMutation(recordID, "submitted", reviewer, "ready for review"); err != nil {
		return domain.Record{}, err
	}
	return record, nil
}

func (s *Service) ConfirmReview(recordID, reviewer, note string) (domain.Record, error) {
	if reviewer == "" {
		return domain.Record{}, errors.New("reviewer is required")
	}
	record, err := s.store.UpdateRecord(recordID, func(target *domain.Record) error {
		if target.Status != domain.StatusReview {
			return fmt.Errorf("record status %s cannot be confirmed", target.Status)
		}
		target.Status = domain.StatusConfirmed
		target.Reviewer = reviewer
		target.UpdatedAt = s.clock.Now()
		target.Version++
		return nil
	})
	if err != nil {
		return domain.Record{}, err
	}
	workflow, err := s.store.GetWorkflowForRecord(recordID)
	if err != nil {
		return domain.Record{}, err
	}
	if _, err := s.store.AdvanceWorkflow(workflow.ID, domain.StageConfirm, note); err != nil {
		return domain.Record{}, err
	}
	if err := s.auditMutation(recordID, "review-confirmed", reviewer, note); err != nil {
		return domain.Record{}, err
	}
	return record, nil
}

func (s *Service) RejectReview(recordID, reviewer, note string) (domain.Record, error) {
	if note == "" {
		return domain.Record{}, errors.New("rejection note is required")
	}
	record, err := s.store.UpdateRecord(recordID, func(target *domain.Record) error {
		if target.Status != domain.StatusReview {
			return errors.New("only review records can be rejected")
		}
		target.Status = domain.StatusRejected
		target.Reviewer = reviewer
		target.UpdatedAt = s.clock.Now()
		target.Version++
		return nil
	})
	if err != nil {
		return domain.Record{}, err
	}
	workflow, err := s.store.GetWorkflowForRecord(recordID)
	if err == nil {
		_, _ = s.store.AdvanceWorkflow(workflow.ID, domain.StageCapture, note)
	}
	if err := s.auditMutation(recordID, "review-rejected", reviewer, note); err != nil {
		return domain.Record{}, err
	}
	return record, nil
}

func (s *Service) ArchiveStory(recordID, actor string) (domain.Record, error) {
	record, err := s.store.UpdateRecord(recordID, func(target *domain.Record) error {
		if target.Status != domain.StatusConfirmed {
			return errors.New("only confirmed records can be archived")
		}
		target.Status = domain.StatusArchived
		target.UpdatedAt = s.clock.Now()
		target.Version++
		return nil
	})
	if err != nil {
		return domain.Record{}, err
	}
	workflow, err := s.store.GetWorkflowForRecord(recordID)
	if err != nil {
		return domain.Record{}, err
	}
	if _, err := s.store.AdvanceWorkflow(workflow.ID, domain.StageArchive, "archived by "+actor); err != nil {
		return domain.Record{}, err
	}
	if err := s.auditMutation(recordID, "archived", actor, "story archived"); err != nil {
		return domain.Record{}, err
	}
	return record, nil
}

func (s *Service) UpdateNarrative(recordID, narrative, actor string) (domain.Record, error) {
	record, err := s.store.UpdateRecord(recordID, func(target *domain.Record) error { return target.UpdateNarrative(narrative, s.clock.Now()) })
	if err != nil {
		return domain.Record{}, err
	}
	if err := s.auditMutation(recordID, "narrative-updated", actor, "story text changed"); err != nil {
		return domain.Record{}, err
	}
	return record, nil
}

func (s *Service) UpdateAmount(recordID string, amount int64, actor string) (domain.Record, error) {
	updated, err := s.store.UpdateRecord(recordID, func(target *domain.Record) error {
		return target.UpdateAmount(amount, s.clock.Now())
	})
	if err != nil {
		return domain.Record{}, err
	}
	if err := s.auditMutation(recordID, "amount-updated", actor, fmt.Sprintf("amount=%d", amount)); err != nil {
		return domain.Record{}, err
	}
	return updated, nil
}

func (s *Service) AddTag(recordID, tag, actor string) (domain.Record, error) {
	record, err := s.store.UpdateRecord(recordID, func(target *domain.Record) error { return target.AddTag(tag, s.clock.Now()) })
	if err != nil {
		return domain.Record{}, err
	}
	if err := s.auditMutation(recordID, "tag-added", actor, tag); err != nil {
		return domain.Record{}, err
	}
	return record, nil
}

func (s *Service) RemoveTag(recordID, tag, actor string) (domain.Record, error) {
	var removed bool
	record, err := s.store.UpdateRecord(recordID, func(target *domain.Record) error { removed = target.RemoveTag(tag, s.clock.Now()); return nil })
	if err != nil {
		return domain.Record{}, err
	}
	if removed {
		if err := s.auditMutation(recordID, "tag-removed", actor, tag); err != nil {
			return domain.Record{}, err
		}
	}
	return record, nil
}

func (s *Service) AddWorkflowNote(recordID, note, actor string) error {
	workflow, err := s.store.GetWorkflowForRecord(recordID)
	if err != nil {
		return err
	}
	if note == "" {
		return errors.New("workflow note is required")
	}
	updated, err := s.store.AdvanceWorkflow(workflow.ID, workflow.Stage, note)
	if err != nil {
		if workflow.Stage == domain.StageArchive {
			return s.auditMutation(recordID, "workflow-note", actor, note)
		}
		return err
	}
	if updated.LatestNote() != note {
		return errors.New("workflow note was not recorded")
	}
	return s.auditMutation(recordID, "workflow-note", actor, note)
}
