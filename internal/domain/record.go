package domain

import (
	"errors"
	"fmt"
	"strings"
)

type RecordStatus string

const (
	StatusDraft     RecordStatus = "draft"
	StatusReview    RecordStatus = "review"
	StatusConfirmed RecordStatus = "confirmed"
	StatusArchived  RecordStatus = "archived"
	StatusRejected  RecordStatus = "rejected"
)

type Record struct {
	ID        string       `json:"id"`
	Title     string       `json:"title"`
	Narrative string       `json:"narrative"`
	Amount    int64        `json:"amount"`
	Currency  string       `json:"currency"`
	Status    RecordStatus `json:"status"`
	Author    string       `json:"author"`
	Reviewer  string       `json:"reviewer"`
	CreatedAt string       `json:"created_at"`
	UpdatedAt string       `json:"updated_at"`
	Version   int          `json:"version"`
	Tags      []string     `json:"tags"`
}

func NewRecord(id, title, narrative, author string, amount int64, createdAt string) Record {
	return Record{ID: id, Title: strings.TrimSpace(title), Narrative: strings.TrimSpace(narrative), Amount: amount, Currency: "CNY", Status: StatusDraft, Author: strings.TrimSpace(author), CreatedAt: createdAt, UpdatedAt: createdAt, Version: 1, Tags: []string{}}
}

func (r Record) Validate() error {
	if strings.TrimSpace(r.ID) == "" {
		return errors.New("record id is required")
	}
	if len([]rune(strings.TrimSpace(r.Title))) < 2 {
		return errors.New("record title must contain at least two characters")
	}
	if len([]rune(strings.TrimSpace(r.Narrative))) < 8 {
		return errors.New("record narrative must contain at least eight characters")
	}
	if r.Amount < 0 {
		return errors.New("record amount cannot be negative")
	}
	if r.Amount > 1000000000 {
		return errors.New("record amount exceeds limit")
	}
	if strings.TrimSpace(r.Currency) == "" {
		return errors.New("record currency is required")
	}
	if strings.TrimSpace(r.Author) == "" {
		return errors.New("record author is required")
	}
	if strings.TrimSpace(r.CreatedAt) == "" || strings.TrimSpace(r.UpdatedAt) == "" {
		return errors.New("record timestamps are required")
	}
	if r.Version < 1 {
		return errors.New("record version must be positive")
	}
	if !r.Status.Valid() {
		return fmt.Errorf("invalid record status %q", r.Status)
	}
	return nil
}

func (s RecordStatus) Valid() bool {
	switch s {
	case StatusDraft, StatusReview, StatusConfirmed, StatusArchived, StatusRejected:
		return true
	default:
		return false
	}
}

func (r Record) IsEditable() bool {
	return r.Status == StatusDraft || r.Status == StatusReview || r.Status == StatusRejected
}

func (r Record) IsVisible() bool {
	return r.Status != StatusRejected
}

func (r Record) Clone() Record {
	copyTags := append([]string(nil), r.Tags...)
	r.Tags = copyTags
	return r
}

func (r Record) Summary() string {
	return fmt.Sprintf("%s: %s (%d %s, %s)", r.ID, r.Title, r.Amount, r.Currency, r.Status)
}

func (r *Record) UpdateNarrative(narrative, now string) error {
	if !r.IsEditable() {
		return errors.New("record is not editable")
	}
	if strings.TrimSpace(narrative) == "" {
		return errors.New("narrative cannot be empty")
	}
	r.Narrative = strings.TrimSpace(narrative)
	r.UpdatedAt = now
	r.Version++
	return nil
}

func (r *Record) UpdateAmount(amount int64, now string) error {
	if !r.IsEditable() {
		return errors.New("record is not editable")
	}
	if amount < 0 || amount > 1000000000 {
		return errors.New("amount outside allowed range")
	}
	r.Amount = amount
	r.UpdatedAt = now
	r.Version++
	return nil
}

func (r *Record) AddTag(tag string, now string) error {
	clean := strings.TrimSpace(tag)
	if clean == "" {
		return errors.New("tag cannot be empty")
	}
	for _, existing := range r.Tags {
		if strings.EqualFold(existing, clean) {
			return nil
		}
	}
	r.Tags = append(r.Tags, clean)
	r.UpdatedAt = now
	r.Version++
	return nil
}

func (r *Record) RemoveTag(tag, now string) bool {
	clean := strings.TrimSpace(tag)
	for i, existing := range r.Tags {
		if strings.EqualFold(existing, clean) {
			r.Tags = append(r.Tags[:i], r.Tags[i+1:]...)
			r.UpdatedAt = now
			r.Version++
			return true
		}
	}
	return false
}

func (r *Record) AssignReviewer(reviewer, now string) error {
	if strings.TrimSpace(reviewer) == "" {
		return errors.New("reviewer is required")
	}
	if r.Status != StatusDraft && r.Status != StatusReview {
		return errors.New("reviewer can only be assigned before confirmation")
	}
	r.Reviewer = strings.TrimSpace(reviewer)
	r.Status = StatusReview
	r.UpdatedAt = now
	r.Version++
	return nil
}
