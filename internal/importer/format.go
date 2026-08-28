package importer

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func ParseCSV(data string) (Bundle, error) {
	reader := csv.NewReader(strings.NewReader(data))
	header, err := reader.Read()
	if err != nil {
		return Bundle{}, err
	}
	columns := map[string]int{}
	for index, name := range header {
		columns[strings.ToLower(strings.TrimSpace(name))] = index
	}
	for _, required := range []string{"title", "narrative", "author", "amount"} {
		if _, ok := columns[required]; !ok {
			return Bundle{}, fmt.Errorf("missing csv column %s", required)
		}
	}
	payload := Payload{Records: []RecordInput{}}
	for {
		row, readErr := reader.Read()
		if errors.Is(readErr, io.EOF) {
			break
		}
		if readErr != nil {
			return Bundle{}, readErr
		}
		amount, parseErr := strconv.ParseInt(strings.TrimSpace(row[columns["amount"]]), 10, 64)
		if parseErr != nil {
			return Bundle{}, fmt.Errorf("amount: %w", parseErr)
		}
		input := RecordInput{Title: row[columns["title"]], Narrative: row[columns["narrative"]], Author: row[columns["author"]], Amount: amount}
		if tagIndex, ok := columns["tags"]; ok && tagIndex < len(row) {
			input.Tags = strings.Split(row[tagIndex], "|")
		}
		payload.Records = append(payload.Records, input)
	}
	return ParsePayload(payload)
}

func ExplainErrors(err error) string {
	if err == nil {
		return ""
	}
	return strings.TrimSpace(err.Error())
}

func CountInputs(payload Payload) (int, int) {
	attachments := 0
	for _, record := range payload.Records {
		attachments += len(record.Attachments)
	}
	return len(payload.Records), attachments
}
