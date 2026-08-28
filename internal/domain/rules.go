package domain

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

type ReviewDecision string

const (
	DecisionApprove ReviewDecision = "approve"
	DecisionReject  ReviewDecision = "reject"
)

type AmountBand string

const (
	AmountSmall  AmountBand = "small"
	AmountMedium AmountBand = "medium"
	AmountLarge  AmountBand = "large"
)

func ValidateTransition(status RecordStatus, decision ReviewDecision) error {
	switch decision {
	case DecisionApprove:
		if status != StatusReview {
			return fmt.Errorf("status %s cannot be approved", status)
		}
	case DecisionReject:
		if status != StatusReview {
			return fmt.Errorf("status %s cannot be rejected", status)
		}
	default:
		return errors.New("unknown review decision")
	}
	return nil
}

func ApplyDecision(record *Record, decision ReviewDecision, reviewer, note, now string) error {
	if record == nil {
		return errors.New("record is nil")
	}
	if err := ValidateTransition(record.Status, decision); err != nil {
		return err
	}
	if strings.TrimSpace(reviewer) == "" {
		return errors.New("reviewer is required")
	}
	if decision == DecisionReject && strings.TrimSpace(note) == "" {
		return errors.New("rejection note is required")
	}
	record.Reviewer = strings.TrimSpace(reviewer)
	if decision == DecisionApprove {
		record.Status = StatusConfirmed
	} else {
		record.Status = StatusRejected
	}
	record.UpdatedAt = now
	record.Version++
	return nil
}

func ClassifyAmount(amount int64) AmountBand {
	if amount < 1000 {
		return AmountSmall
	}
	if amount < 10000 {
		return AmountMedium
	}
	return AmountLarge
}

func SortRecords(records []Record, descending bool) []Record {
	result := make([]Record, len(records))
	for i, record := range records {
		result[i] = record.Clone()
	}
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Amount == result[j].Amount {
			return result[i].ID < result[j].ID
		}
		if descending {
			return result[i].Amount > result[j].Amount
		}
		return result[i].Amount < result[j].Amount
	})
	return result
}

func MergeTags(left, right []string) []string {
	combined := append(append([]string{}, left...), right...)
	return NormalizeTags(combined)
}

func CanonicalTitle(title string) string {
	words := strings.Fields(strings.TrimSpace(title))
	return strings.Join(words, " ")
}

func NormalizeActor(actor string) string {
	clean := strings.TrimSpace(actor)
	if clean == "" {
		return "system"
	}
	return clean
}

func CanArchive(record Record) bool {
	return record.Status == StatusConfirmed && record.Reviewer != ""
}

func NarrativeWords(narrative string) []string {
	words := strings.Fields(strings.TrimSpace(narrative))
	result := make([]string, 0, len(words))
	for _, word := range words {
		if len([]rune(word)) >= 2 {
			result = append(result, word)
		}
	}
	return result
}

func EnsureTags(record *Record) {
	if record == nil {
		return
	}
	record.Tags = NormalizeTags(record.Tags)
}
