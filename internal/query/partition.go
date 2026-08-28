package query

import (
	"sort"

	"genealogy-story-organizer/internal/domain"
)

type Partition struct {
	Visible  []domain.Record
	Rejected []domain.Record
	Editable []domain.Record
	Archived []domain.Record
}

func PartitionRecords(records []domain.Record) Partition {
	partition := Partition{Visible: []domain.Record{}, Rejected: []domain.Record{}, Editable: []domain.Record{}, Archived: []domain.Record{}}
	for _, record := range records {
		clone := record.Clone()
		if record.IsVisible() {
			partition.Visible = append(partition.Visible, clone)
		} else {
			partition.Rejected = append(partition.Rejected, clone)
		}
		if record.IsEditable() {
			partition.Editable = append(partition.Editable, clone)
		}
		if record.Status == domain.StatusArchived {
			partition.Archived = append(partition.Archived, clone)
		}
	}
	sort.SliceStable(partition.Visible, func(i, j int) bool { return partition.Visible[i].ID < partition.Visible[j].ID })
	sort.SliceStable(partition.Rejected, func(i, j int) bool { return partition.Rejected[i].ID < partition.Rejected[j].ID })
	return partition
}

func (p Partition) CountVisible() int { return len(p.Visible) }

func (p Partition) CountRejected() int { return len(p.Rejected) }

func (p Partition) HasArchived() bool { return len(p.Archived) > 0 }

func (p Partition) AmountTotal() int64 {
	var total int64
	for _, record := range p.Visible {
		total += record.Amount
	}
	return total
}
