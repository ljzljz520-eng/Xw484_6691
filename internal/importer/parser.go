package importer

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"genealogy-story-organizer/internal/domain"
)

type AttachmentInput struct {
	Name      string `json:"name"`
	MediaType string `json:"media_type"`
	Checksum  string `json:"checksum"`
	Size      int64  `json:"size"`
}

type RecordInput struct {
	Title       string            `json:"title"`
	Narrative   string            `json:"narrative"`
	Author      string            `json:"author"`
	Amount      int64             `json:"amount"`
	Tags        []string          `json:"tags"`
	Attachments []AttachmentInput `json:"attachments"`
}

type Payload struct {
	Records []RecordInput `json:"records"`
}

type Bundle struct{ Records []RecordInput }

func ParseJSON(data []byte) (Bundle, error) {
	if len(data) == 0 {
		return Bundle{}, errors.New("payload is empty")
	}
	var payload Payload
	if err := json.Unmarshal(data, &payload); err != nil {
		return Bundle{}, fmt.Errorf("decode payload: %w", err)
	}
	return ParsePayload(payload)
}

func ParsePayload(payload Payload) (Bundle, error) {
	if len(payload.Records) == 0 {
		return Bundle{}, errors.New("payload contains no records")
	}
	result := Bundle{Records: make([]RecordInput, len(payload.Records))}
	for i, input := range payload.Records {
		input.Title = strings.TrimSpace(input.Title)
		input.Narrative = strings.TrimSpace(input.Narrative)
		input.Author = strings.TrimSpace(input.Author)
		input.Tags = domain.NormalizeTags(input.Tags)
		for j := range input.Attachments {
			input.Attachments[j].Name = strings.TrimSpace(input.Attachments[j].Name)
			input.Attachments[j].MediaType = strings.ToLower(strings.TrimSpace(input.Attachments[j].MediaType))
			input.Attachments[j].Checksum = strings.TrimSpace(input.Attachments[j].Checksum)
		}
		result.Records[i] = input
	}
	return result, nil
}

func EncodePayload(payload Payload) ([]byte, error) { return json.MarshalIndent(payload, "", "  ") }
