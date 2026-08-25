package application

import (
	"fmt"

	"genealogy-story-organizer/internal/domain"
	"genealogy-story-organizer/internal/importer"
	"genealogy-story-organizer/internal/query"
)

type ImportResult struct {
	Records     []domain.Record
	Attachments []domain.Attachment
	Report      query.Report
}

func (s *Service) ImportPayload(payload importer.Payload) (ImportResult, error) {
	parsed, err := importer.ParsePayload(payload)
	if err != nil {
		return ImportResult{}, err
	}
	if err := importer.ValidateBundle(parsed); err != nil {
		return ImportResult{}, err
	}
	result := ImportResult{Records: []domain.Record{}, Attachments: []domain.Attachment{}}
	for _, input := range parsed.Records {
		record, err := s.CreateStory(input.Title, input.Narrative, input.Author, input.Amount, input.Tags)
		if err != nil {
			return ImportResult{}, err
		}
		result.Records = append(result.Records, record)
		for _, item := range input.Attachments {
			attachment, err := s.Attach(record.ID, item.Name, item.MediaType, item.Checksum, item.Size)
			if err != nil {
				return ImportResult{}, err
			}
			result.Attachments = append(result.Attachments, attachment)
		}
		workflow, workflowErr := s.store.GetWorkflowForRecord(record.ID)
		if workflowErr != nil {
			return ImportResult{}, workflowErr
		}
		if _, workflowErr = s.store.AdvanceWorkflow(workflow.ID, domain.StageImported, "imported payload"); workflowErr != nil {
			return ImportResult{}, workflowErr
		}
	}
	result.Report = query.BuildReport(result.Records, result.Attachments)
	return result, nil
}

func (s *Service) Explain(recordID string) (string, error) {
	record, err := s.store.GetRecord(recordID)
	if err != nil {
		return "", err
	}
	audit, err := s.store.ListAudit(recordID)
	if err != nil {
		return "", err
	}
	workflow, err := s.store.GetWorkflowForRecord(recordID)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s | stage=%s | events=%d", record.Summary(), domain.StageLabel(workflow.Stage), len(audit)), nil
}
