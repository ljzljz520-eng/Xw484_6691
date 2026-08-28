package importer

import (
	"errors"
	"fmt"
	"strings"

	"genealogy-story-organizer/internal/domain"
)

func ValidateBundle(bundle Bundle) error {
	if len(bundle.Records) == 0 {
		return errors.New("bundle is empty")
	}
	seenTitles := map[string]bool{}
	for index, input := range bundle.Records {
		if err := validateRecordInput(input); err != nil {
			return fmt.Errorf("record %d: %w", index+1, err)
		}
		key := strings.ToLower(input.Title)
		if seenTitles[key] {
			return fmt.Errorf("record %q is duplicated", input.Title)
		}
		seenTitles[key] = true
		for attachmentIndex, item := range input.Attachments {
			if err := validateAttachmentInput(item); err != nil {
				return fmt.Errorf("record %d attachment %d: %w", index+1, attachmentIndex+1, err)
			}
		}
	}
	return nil
}

func validateRecordInput(input RecordInput) error {
	if len([]rune(input.Title)) < 2 {
		return errors.New("title is too short")
	}
	if len([]rune(input.Narrative)) < 8 {
		return errors.New("narrative is too short")
	}
	if input.Amount < 0 || input.Amount > 1000000000 {
		return errors.New("amount is outside range")
	}
	if input.Author == "" {
		return errors.New("author is required")
	}
	return nil
}

func validateAttachmentInput(input AttachmentInput) error {
	attachment := domain.NewAttachment("import", "record", input.Name, input.MediaType, input.Checksum, "2000-01-01T00:00:00Z", input.Size)
	return attachment.Validate()
}

func ValidateSingle(input RecordInput) []string {
	errorsFound := []string{}
	if input.Title == "" {
		errorsFound = append(errorsFound, "title")
	}
	if input.Narrative == "" {
		errorsFound = append(errorsFound, "narrative")
	}
	if input.Author == "" {
		errorsFound = append(errorsFound, "author")
	}
	if input.Amount < 0 {
		errorsFound = append(errorsFound, "amount")
	}
	return errorsFound
}
