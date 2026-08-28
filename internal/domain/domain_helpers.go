package domain

import "strings"

func NormalizeTags(tags []string) []string {
	seen := map[string]bool{}
	result := make([]string, 0, len(tags))
	for _, raw := range tags {
		clean := strings.TrimSpace(raw)
		if clean == "" {
			continue
		}
		key := strings.ToLower(clean)
		if seen[key] {
			continue
		}
		seen[key] = true
		result = append(result, clean)
	}
	return result
}

func StatusLabel(status RecordStatus) string {
	switch status {
	case StatusDraft:
		return "草稿"
	case StatusReview:
		return "审核中"
	case StatusConfirmed:
		return "已确认"
	case StatusArchived:
		return "已归档"
	case StatusRejected:
		return "已驳回"
	default:
		return "未知"
	}
}

func StageLabel(stage WorkflowStage) string {
	switch stage {
	case StageCapture:
		return "登记"
	case StageReview:
		return "审核"
	case StageConfirm:
		return "确认"
	case StageArchive:
		return "归档"
	case StageImported:
		return "导入"
	default:
		return "未知"
	}
}
