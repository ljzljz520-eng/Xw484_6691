package domain

import (
	"sort"
	"strings"
)

type Index struct {
	ByID     map[string]Record
	ByTag    map[string][]string
	ByStatus map[RecordStatus][]string
	ByAuthor map[string][]string
}

func BuildIndex(records []Record) Index {
	index := Index{ByID: map[string]Record{}, ByTag: map[string][]string{}, ByStatus: map[RecordStatus][]string{}, ByAuthor: map[string][]string{}}
	for _, record := range records {
		index.ByID[record.ID] = record.Clone()
		index.ByStatus[record.Status] = append(index.ByStatus[record.Status], record.ID)
		author := strings.ToLower(strings.TrimSpace(record.Author))
		index.ByAuthor[author] = append(index.ByAuthor[author], record.ID)
		for _, tag := range NormalizeTags(record.Tags) {
			key := strings.ToLower(tag)
			index.ByTag[key] = append(index.ByTag[key], record.ID)
		}
	}
	index.sortLists()
	return index
}

func (i *Index) sortLists() {
	for key, ids := range i.ByTag {
		sort.Strings(ids)
		i.ByTag[key] = ids
	}
	for key, ids := range i.ByStatus {
		sort.Strings(ids)
		i.ByStatus[key] = ids
	}
	for key, ids := range i.ByAuthor {
		sort.Strings(ids)
		i.ByAuthor[key] = ids
	}
}

func (i Index) IDsForTag(tag string) []string {
	return append([]string(nil), i.ByTag[strings.ToLower(strings.TrimSpace(tag))]...)
}

func (i Index) IDsForStatus(status RecordStatus) []string {
	return append([]string(nil), i.ByStatus[status]...)
}

func (i Index) IDsForAuthor(author string) []string {
	return append([]string(nil), i.ByAuthor[strings.ToLower(strings.TrimSpace(author))]...)
}

func (i Index) Lookup(id string) (Record, bool) { record, ok := i.ByID[id]; return record.Clone(), ok }

func (i Index) Size() int { return len(i.ByID) }

func (i Index) Empty() bool { return i.Size() == 0 }
