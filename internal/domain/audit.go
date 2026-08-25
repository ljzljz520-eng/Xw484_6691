package domain

import (
	"errors"
	"strings"
)

type AuditEvent struct {
	ID        string `json:"id"`
	RecordID  string `json:"record_id"`
	Action    string `json:"action"`
	Actor     string `json:"actor"`
	Note      string `json:"note"`
	CreatedAt string `json:"created_at"`
	Sequence  int    `json:"sequence"`
}

func NewAuditEvent(id, recordID, action, actor, note, at string, sequence int) AuditEvent {
	return AuditEvent{ID: id, RecordID: recordID, Action: strings.TrimSpace(action), Actor: strings.TrimSpace(actor), Note: strings.TrimSpace(note), CreatedAt: at, Sequence: sequence}
}

func (a AuditEvent) Validate() error {
	if strings.TrimSpace(a.ID) == "" || strings.TrimSpace(a.RecordID) == "" {
		return errors.New("audit identifiers are required")
	}
	if strings.TrimSpace(a.Action) == "" || strings.TrimSpace(a.Actor) == "" {
		return errors.New("audit action and actor are required")
	}
	if strings.TrimSpace(a.CreatedAt) == "" {
		return errors.New("audit timestamp is required")
	}
	if a.Sequence < 1 {
		return errors.New("audit sequence must be positive")
	}
	return nil
}

func (a AuditEvent) IsMutation() bool {
	switch strings.ToLower(a.Action) {
	case "amount-updated", "narrative-updated", "tag-added", "tag-removed":
		return true
	default:
		return false
	}
}

func (a AuditEvent) Label() string {
	if a.Note == "" {
		return a.Action + " by " + a.Actor
	}
	return a.Action + " by " + a.Actor + ": " + a.Note
}

func SortAudit(events []AuditEvent) []AuditEvent {
	result := append([]AuditEvent(nil), events...)
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j].Sequence < result[i].Sequence {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	return result
}
