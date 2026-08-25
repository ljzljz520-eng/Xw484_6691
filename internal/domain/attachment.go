package domain

import (
	"errors"
	"strings"
)

type Attachment struct {
	ID        string `json:"id"`
	RecordID  string `json:"record_id"`
	Name      string `json:"name"`
	MediaType string `json:"media_type"`
	Checksum  string `json:"checksum"`
	Size      int64  `json:"size"`
	CreatedAt string `json:"created_at"`
}

func NewAttachment(id, recordID, name, mediaType, checksum, at string, size int64) Attachment {
	return Attachment{ID: id, RecordID: recordID, Name: strings.TrimSpace(name), MediaType: strings.TrimSpace(mediaType), Checksum: strings.TrimSpace(checksum), Size: size, CreatedAt: at}
}

func (a Attachment) Validate() error {
	if strings.TrimSpace(a.ID) == "" || strings.TrimSpace(a.RecordID) == "" {
		return errors.New("attachment identifiers are required")
	}
	if strings.TrimSpace(a.Name) == "" || strings.TrimSpace(a.MediaType) == "" {
		return errors.New("attachment name and media type are required")
	}
	if !SupportedMediaType(a.MediaType) {
		return errors.New("attachment media type is unsupported")
	}
	if strings.TrimSpace(a.Checksum) == "" {
		return errors.New("attachment checksum is required")
	}
	if a.Size < 0 || a.Size > 50*1024*1024 {
		return errors.New("attachment size is out of range")
	}
	if strings.TrimSpace(a.CreatedAt) == "" {
		return errors.New("attachment timestamp is required")
	}
	return nil
}

func SupportedMediaType(mediaType string) bool {
	switch strings.ToLower(strings.TrimSpace(mediaType)) {
	case "text/plain", "text/markdown", "application/pdf", "image/jpeg", "image/png":
		return true
	default:
		return false
	}
}

func (a Attachment) IsImage() bool {
	return strings.HasPrefix(strings.ToLower(a.MediaType), "image/")
}

func (a Attachment) DisplayName() string {
	if strings.TrimSpace(a.Name) == "" {
		return a.ID
	}
	return a.Name
}
