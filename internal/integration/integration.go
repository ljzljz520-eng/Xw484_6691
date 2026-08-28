package integration

import (
	"fmt"

	"genealogy-story-organizer/internal/application"
	"genealogy-story-organizer/internal/domain"
)

type Snapshot struct {
	Record   domain.Record
	Workflow domain.Workflow
	Audit    []domain.AuditEvent
}

func BuildSnapshot(service *application.Service, recordID string) (Snapshot, error) {
	record, err := service.GetStory(recordID)
	if err != nil {
		return Snapshot{}, err
	}
	workflow, err := service.Database().GetWorkflowForRecord(recordID)
	if err != nil {
		return Snapshot{}, err
	}
	audit, err := service.Database().ListAudit(recordID)
	if err != nil {
		return Snapshot{}, err
	}
	return Snapshot{Record: record, Workflow: workflow, Audit: audit}, nil
}

func (s Snapshot) Complete() bool {
	return s.Record.Status == domain.StatusArchived && s.Workflow.IsComplete()
}

func (s Snapshot) Narrative() string {
	return fmt.Sprintf("%s/%s/%d", s.Record.Title, s.Workflow.Stage, len(s.Audit))
}
