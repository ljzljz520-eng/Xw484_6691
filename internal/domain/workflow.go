package domain

import (
	"errors"
	"strings"
)

type WorkflowStage string

const (
	StageCapture  WorkflowStage = "capture"
	StageReview   WorkflowStage = "review"
	StageConfirm  WorkflowStage = "confirm"
	StageArchive  WorkflowStage = "archive"
	StageImported WorkflowStage = "imported"
)

type Workflow struct {
	ID        string        `json:"id"`
	RecordID  string        `json:"record_id"`
	Stage     WorkflowStage `json:"stage"`
	Owner     string        `json:"owner"`
	UpdatedAt string        `json:"updated_at"`
	Notes     []string      `json:"notes"`
	Revision  int           `json:"revision"`
}

func NewWorkflow(id, recordID, owner, now string) Workflow {
	return Workflow{ID: id, RecordID: recordID, Stage: StageCapture, Owner: owner, UpdatedAt: now, Notes: []string{}, Revision: 1}
}

func (w Workflow) Validate() error {
	if strings.TrimSpace(w.ID) == "" || strings.TrimSpace(w.RecordID) == "" {
		return errors.New("workflow identifiers are required")
	}
	if !w.Stage.Valid() {
		return errors.New("workflow stage is invalid")
	}
	if strings.TrimSpace(w.Owner) == "" || strings.TrimSpace(w.UpdatedAt) == "" {
		return errors.New("workflow owner and timestamp are required")
	}
	if w.Revision < 1 {
		return errors.New("workflow revision must be positive")
	}
	return nil
}

func (s WorkflowStage) Valid() bool {
	switch s {
	case StageCapture, StageReview, StageConfirm, StageArchive, StageImported:
		return true
	default:
		return false
	}
}

func (w *Workflow) Advance(next WorkflowStage, note, now string) error {
	if !next.Valid() {
		return errors.New("unknown workflow stage")
	}
	if !w.canAdvance(next) {
		return errors.New("workflow stage transition is not allowed")
	}
	w.Stage = next
	w.UpdatedAt = now
	w.Revision++
	if clean := strings.TrimSpace(note); clean != "" {
		w.Notes = append(w.Notes, clean)
	}
	return nil
}

func (w Workflow) canAdvance(next WorkflowStage) bool {
	switch w.Stage {
	case StageCapture:
		return next == StageReview || next == StageImported
	case StageReview:
		return next == StageConfirm || next == StageCapture
	case StageConfirm:
		return next == StageArchive || next == StageReview
	case StageImported:
		return next == StageReview || next == StageArchive
	case StageArchive:
		return false
	default:
		return false
	}
}

func (w Workflow) IsComplete() bool { return w.Stage == StageArchive }

func (w Workflow) LatestNote() string {
	if len(w.Notes) == 0 {
		return ""
	}
	return w.Notes[len(w.Notes)-1]
}
