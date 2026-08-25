package application

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"genealogy-story-organizer/internal/domain"
	"genealogy-story-organizer/internal/query"
)

type AuditDigest struct {
	RecordID  string
	Events    int
	Mutations int
	Actors    []string
	Latest    string
}

func (s *Service) AuditDigest(recordID string) (AuditDigest, error) {
	if strings.TrimSpace(recordID) == "" {
		return AuditDigest{}, errors.New("record id is required")
	}
	events, err := s.store.ListAudit(recordID)
	if err != nil {
		return AuditDigest{}, err
	}
	actors := map[string]bool{}
	mutations := 0
	for _, event := range events {
		actors[event.Actor] = true
		if event.IsMutation() {
			mutations++
		}
	}
	actorList := make([]string, 0, len(actors))
	for actor := range actors {
		actorList = append(actorList, actor)
	}
	sort.Strings(actorList)
	latest := ""
	if len(events) > 0 {
		latest = events[len(events)-1].Action
	}
	return AuditDigest{RecordID: recordID, Events: len(events), Mutations: mutations, Actors: actorList, Latest: latest}, nil
}

func (s *Service) ChangeHistory(recordID string) (string, error) {
	digest, err := s.AuditDigest(recordID)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s events=%d mutations=%d latest=%s actors=%s", digest.RecordID, digest.Events, digest.Mutations, digest.Latest, strings.Join(digest.Actors, ",")), nil
}

func (s *Service) SearchAndSummarize(filter query.Filter) ([]string, error) {
	records, err := s.store.ListRecords()
	if err != nil {
		return nil, err
	}
	found := query.Search(records, filter)
	rows := query.ToExportRows(found)
	result := make([]string, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.ID+" "+row.Title+" "+domain.StatusLabel(domain.RecordStatus(row.Status)))
	}
	return result, nil
}

func (s *Service) ValidateStorage() error {
	problems, err := s.store.ValidateReferences()
	if err != nil {
		return err
	}
	if len(problems) > 0 {
		return fmt.Errorf("storage references invalid: %s", strings.Join(problems, ","))
	}
	return nil
}
